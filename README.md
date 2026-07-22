# Mini Kubernetes - Phase 1

## Implementation: Go + REST + gRPC

This repo implements the Phase 1 spec in Go:

- **CLI** (`main.go`): entry point that sends HTTP requests
- **Control Plane** (`pkg/controlplane`): REST server that validates and forwards requests to worker via gRPC
- **Worker** (`pkg/worker`): gRPC server that calls Docker

**Why this split?**

- REST for the user-facing API (simple, debuggable)
- gRPC for internal communication (typed, efficient)

### Build

```bash
go build -o mini-kube
```

### Run

Terminal 1 (worker):

```bash
./mini-kube worker
```

Terminal 2 (control plane):

```bash
./mini-kube control-plane
```

Terminal 3 (CLI):

```bash
./mini-kube deploy nginx
```

---

## Overview

The goal of this project is **not** to recreate Kubernetes.

The goal is to understand **why Kubernetes exists** by building a simplified version from scratch.

Over multiple phases, we will build our own small container orchestrator and gradually add features similar to Kubernetes.

---

# Phase 1 - Remote Container Deployment

## Objective

Build the smallest possible container orchestrator.

Instead of running Docker manually:

```bash
docker run nginx
```

a user should be able to run:

```bash
mini-kube deploy nginx
```

Our system should handle everything else.

---

# Problem Statement

Imagine you have many machines running Docker.

Instead of logging into every machine and executing Docker commands manually, we want a central system that accepts deployment requests and instructs workers to run containers.

For Phase 1, we only have **one worker**, but the architecture should allow us to expand later.

---

# Functional Requirements

The system should contain three logical components.

## 1. Control Plane

The Control Plane is the brain of the system.

Responsibilities:

- Accept deployment requests.
- Validate requests.
- Forward deployment requests to a worker.
- Return success or failure to the user.

The Control Plane **must not** run Docker directly.

---

## 2. Worker

The Worker is responsible for executing workloads.

Responsibilities:

- Wait for deployment requests.
- Receive a Docker image name.
- Execute Docker.
- Report success or failure back to the Control Plane.

Example:

```text
Receive:

{
    "image": "nginx"
}

↓

Execute:

docker run nginx

↓

Return success
```

---

## 3. Communication Layer

The Control Plane and Worker must communicate over the network.

You may choose:

- REST
- gRPC
- Raw TCP sockets

Be prepared to explain why you chose your approach.

---

# End-to-End Flow

The completed system should work like this:

```text
User

↓

mini-kube deploy nginx

↓

Control Plane

↓

Worker

↓

Docker

↓

Running Container

↓

Worker reports success

↓

Control Plane reports success

↓

User sees deployment successful
```

---

# Expected Architecture

Example architecture:

```text
                User
                  │
                  ▼
       mini-kube deploy nginx
                  │
                  ▼
        +--------------------+
        |   Control Plane    |
        +--------------------+
                  │
          Deployment Request
                  │
                  ▼
        +--------------------+
        |       Worker       |
        +--------------------+
                  │
          docker run nginx
                  │
                  ▼
        +--------------------+
        |      Docker        |
        +--------------------+
                  │
                  ▼
         Running Container
```

---

# Project Scope

For Phase 1, support only:

- Deploying a Docker image

Nothing else.

---

# Out of Scope

Do **not** implement:

- Multiple workers
- Scheduling
- Scaling
- Desired state
- Health checks
- Heartbeats
- Service discovery
- Load balancing
- Rolling updates
- Cluster persistence
- High availability

Those will be implemented in future phases.

---

# Suggested API

Example request:

```
POST /deploy
```

Body:

```json
{
  "image": "nginx"
}
```

Response:

```json
{
  "status": "success"
}
```

The exact API design is up to you.

---

# Success Criteria

Phase 1 is complete if the following flow works:

```text
User

↓

Deploy request

↓

Control Plane receives request

↓

Worker receives request

↓

Worker starts Docker container

↓

Container is running

↓

Success returned to user
```

---

# Deliverables

Each team member should build their own implementation independently.

Each implementation should include:

- Control Plane
- Worker
- Communication layer
- Docker integration

No code sharing until presentations.

---

# Presentation Requirements

Each person should present:

## 1. Architecture

Draw the system.

Explain every component.

---

## 2. Request Flow

Walk through exactly what happens after:

```bash
mini-kube deploy nginx
```

Explain every step until the container is running.

---

## 3. API Design

Explain:

- Endpoints
- Request format
- Response format

---

## 4. Design Decisions

Be prepared to answer:

- Why this architecture?
- Why this communication protocol?
- Why separate the Control Plane from the Worker?
- How is Docker executed?
- How are failures handled?

---

## 5. Live Demo

Demonstrate:

```bash
mini-kube deploy nginx
```

Show that:

- The request reaches the Control Plane.
- The Worker receives it.
- Docker starts the container.
- The container is successfully running.

---

# Learning Objectives

By completing Phase 1, everyone should understand:

- Client-server communication
- REST or gRPC
- Docker fundamentals
- Linux processes
- Remote execution
- Basic distributed system architecture

---

# Important Rule

The goal is **not** to produce identical implementations.

Each person should design the system independently.

During the review session, compare:

- Architecture
- API design
- Folder structure
- Communication protocol
- Error handling
- Code organization

Discuss the trade-offs and decide what ideas should be carried forward into the team's shared implementation.

---

# End Goal of Phase 1

Build a minimal distributed system where:

- A user requests deployment.
- A Control Plane receives the request.
- A Worker executes the workload.
- Docker starts the container.
- The user receives confirmation.

Once this works, we have built the foundation that every later feature (scheduling, health monitoring, scaling, reconciliation, etc.) will build upon.
