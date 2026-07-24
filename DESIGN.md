# Mini Kubernetes – Design Decisions (Phase 1)

## 1. Why This Architecture?

The system follows a **three-tier architecture**: CLI → Control Plane → Worker → Docker.

This mirrors the real Kubernetes model (kubectl → API Server → kubelet) at the smallest possible scale. Even though Phase 1 only has one worker, the separation means future phases can add workers, a scheduler, and a reconciliation loop without rewriting the core.

The key insight is **separation of concerns**:

| Layer | Responsibility | What it does NOT know |
|---|---|---|
| CLI | User input/output | How containers are started |
| Control Plane | Validation, routing, orchestration | Which runtime is used (Docker, containerd, etc.) |
| Worker | Container lifecycle execution | Whether there are other workers |
| DockerRunner | Shell out to `docker run` | HTTP, routing, validation |

Each layer can evolve independently. Replacing Docker with containerd, or REST with gRPC, is a change in one file.

---

## 2. Why REST over HTTP?

Phase 1 has exactly **one endpoint** (`POST /deploy` → `POST /run`). Three options were considered:

| Option | Pros | Cons |
|---|---|---|
| **REST** | Zero deps, `curl`-testable, human-readable, Go stdlib | No streaming, no bidirectional comms |
| **gRPC** | Strong typing, streaming, efficient binary encoding | Needs protobuf codegen toolchain, harder to debug |
| **Raw TCP** | Maximum control, lowest overhead | Must design framing protocol, error-prone |

**REST wins for Phase 1** because:
- Go's `net/http` is production-quality with zero external dependencies.
- `curl -X POST localhost:8080/deploy -d '{"image":"nginx"}'` is the fastest way to test.
- The overhead of JSON serialization is irrelevant at this scale.
- Swapping to gRPC in Phase 2+ requires only changing transport layer — the `ControlPlane` and `Worker` structs stay the same.

---

## 3. Why Separate the Control Plane from the Worker?

Even in Phase 1 (single worker), this separation provides:

1. **Extensibility**: Adding a second worker means adding a worker URL to the Control Plane — not rewriting it.
2. **Validation boundary**: The Control Plane validates before forwarding. The Worker trusts its input (it only receives pre-validated requests from the Control Plane).
3. **Failure isolation**: If Docker crashes on the Worker, the Control Plane remains up and can report the error gracefully.
4. **Matches production reality**: In a real cluster, the API server and kubelets are always separate processes on separate machines.

---

## 4. How is Docker Executed?

The `DockerRunner` uses `os/exec` to shell out to the `docker` CLI:

```go
cmd := exec.Command("docker", "run", "-d", image)
output, err := cmd.CombinedOutput()
```

**Why `os/exec` instead of the Docker SDK?**

- The Docker SDK (`github.com/docker/docker/client`) adds ~50+ transitive dependencies.
- For Phase 1, we literally run one command: `docker run -d <image>`.
- `os/exec` is stdlib, zero dependencies, and trivially understandable.
- The `Runner` interface means we can swap to the SDK in a later phase.

**Why `-d` (detached mode)?**

- The CLI must return immediately with the container ID.
- Without `-d`, `docker run` blocks until the container exits — unusable for long-running services like nginx.

---

## 5. How Are Failures Handled?

Errors propagate cleanly through every layer:

```text
Docker fails     →  ExecRunner returns error
                 →  Worker sends {"status":"error","error":"..."}
                 →  ControlPlane relays {"status":"error","error":"worker communication failed: ..."}
                 →  CLI prints "❌ Deployment failed!" and exits 1
```

Specific failure cases:

| Failure | Handled By | Response |
|---|---|---|
| Empty image name | Control Plane validation | 400 Bad Request |
| Invalid JSON body | Control Plane / Worker JSON decoder | 400 Bad Request |
| Worker unreachable | Control Plane HTTP client | 502 Bad Gateway |
| Docker not installed | ExecRunner `os/exec` | 500 + error message |
| Image doesn't exist | Docker CLI stderr | 500 + "docker run failed: ..." |
| Wrong HTTP method | Both servers' method check | 405 Method Not Allowed |

---

## 6. Design Patterns Used

### Command Pattern — `DeployRequest`

The deployment intent is captured as a **serializable data object** that flows across process boundaries (CLI → Control Plane → Worker). Each layer can inspect, validate, log, and forward it without coupling to the caller.

### Façade Pattern — `ControlPlane`

The Control Plane exposes **one simple endpoint** to the outside world. Behind that endpoint, it orchestrates validation, worker selection (trivial in Phase 1), HTTP forwarding, and error mapping. The CLI doesn't know workers exist.

### Strategy Pattern — `docker.Runner` Interface

```go
type Runner interface {
    Run(image string) (containerID string, err error)
}
```

The Worker depends on this interface, not on `ExecRunner` directly. This enables:
- **Testing**: Inject a `FakeRunner` that returns a canned container ID.
- **Runtime swapping**: Replace Docker with containerd by implementing a new `Runner`.

### Single Responsibility Principle

Each struct does exactly one thing:

| Struct | Single Responsibility |
|---|---|
| `DeployRequest` / `DeployResponse` | Carry data |
| `ControlPlane` | Accept, validate, forward |
| `Worker` | Receive, execute, report |
| `ExecRunner` | Shell out to Docker |

---

## 7. Class / Struct Summary

```text
┌─────────────────────────────────────────────────────────┐
│                        models                           │
│  DeployRequest   { Image string }                       │
│  DeployResponse  { Status, ContainerID, Error string }  │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    controlplane                         │
│  ControlPlane    { workerURL string }                   │
│    .HandleDeploy(w, r)                                  │
│    .Start(addr)                                         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                       worker                            │
│  Worker          { runner docker.Runner }               │
│    .HandleRun(w, r)                                     │
│    .Start(addr)                                         │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                       docker                            │
│  Runner (interface)                                     │
│    Run(image) → (containerID, error)                    │
│                                                         │
│  ExecRunner (implements Runner)                         │
│    .Run(image) — exec "docker run -d <image>"           │
└─────────────────────────────────────────────────────────┘
```

**Total: 4 structs + 1 interface**, organized across 4 packages.

---

## 8. How to Run

```bash
# Terminal 1 — Start the Worker
go run ./cmd/worker

# Terminal 2 — Start the Control Plane
go run ./cmd/controlplane

# Terminal 3 — Deploy a container
go run ./cmd/cli deploy nginx
```
