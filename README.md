# Mini Kubernetes - Phase 1

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

* Accept deployment requests.
* Validate requests.
* Forward deployment requests to a worker.
* Return success or failure to the user.

The Control Plane **must not** run Docker directly.

---

## 2. Worker

The Worker is responsible for executing workloads.

Responsibilities:

* Wait for deployment requests.
* Receive a Docker image name.
* Execute Docker.
* Report success or failure back to the Control Plane.

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

* REST
* gRPC
* Raw TCP sockets

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

* Deploying a Docker image

Nothing else.

---

# Out of Scope

Do **not** implement:

* Multiple workers
* Scheduling
* Scaling
* Desired state
* Health checks
* Heartbeats
* Service discovery
* Load balancing
* Rolling updates
* Cluster persistence
* High availability

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

* Control Plane
* Worker
* Communication layer
* Docker integration

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

* Endpoints
* Request format
* Response format

---

## 4. Design Decisions

Be prepared to answer:

* Why this architecture?
* Why this communication protocol?
* Why separate the Control Plane from the Worker?
* How is Docker executed?
* How are failures handled?

---

## 5. Live Demo

Demonstrate:

```bash
mini-kube deploy nginx
```

Show that:

* The request reaches the Control Plane.
* The Worker receives it.
* Docker starts the container.
* The container is successfully running.

---

# Learning Objectives

By completing Phase 1, everyone should understand:

* Client-server communication
* REST or gRPC
* Docker fundamentals
* Linux processes
* Remote execution
* Basic distributed system architecture

---

# Important Rule

The goal is **not** to produce identical implementations.

Each person should design the system independently.

During the review session, compare:

* Architecture
* API design
* Folder structure
* Communication protocol
* Error handling
* Code organization

Discuss the trade-offs and decide what ideas should be carried forward into the team's shared implementation.

---

# End Goal of Phase 1

Build a minimal distributed system where:

* A user requests deployment.
* A Control Plane receives the request.
* A Worker executes the workload.
* Docker starts the container.
* The user receives confirmation.

Once this works, we have built the foundation that every later feature (scheduling, health monitoring, scaling, reconciliation, etc.) will build upon.

# Phase 2

## Overview

The goal of this project is not to recreate Kubernetes.

The goal is to understand why Kubernetes exists by building a simplified version from scratch.

In Phase 1, we built a system that can remotely deploy a Docker container through a Control Plane and Worker architecture.

However, the system only performs actions once.

If a container crashes, the system does nothing.

Phase 2 introduces one of the most important concepts in Kubernetes:

**Desired State and Reconciliation.**

The system will now continuously monitor the current state and make changes until the actual state matches the desired state.

---

# Phase 2 - Desired State + Reconciliation Loop

## Objective

Upgrade the system from:

> "Start this container"

to:

> "Keep this container running according to what the user requested."

Instead of:

```bash
docker run nginx
```

the user should be able to define:

```bash
mini-kube create deployment nginx --replicas 3
```

Our system should ensure that 3 nginx containers are always running.

---

# Problem Statement

Imagine you manage multiple applications running inside containers.

A user does not want to manually restart containers whenever they fail.

Instead, they want to declare:

> "I want 3 copies of nginx running."

The system should:

- Remember what the user requested.
- Check what is currently running.
- Detect differences.
- Automatically fix those differences.

For Phase 2, the system still uses one worker, but the architecture should prepare for future scaling.

---

# Functional Requirements

The system should contain the following components.

---

# 1. Deployment Object

The Deployment Object represents what the user wants.

Responsibilities:

- Store application information.
- Store container image information.
- Store desired replica count.
- Represent the desired state of the system.

Example:

```json
{
 "name": "nginx",
 "image": "nginx",
 "replicas": 3
}
```

Meaning:

"The user wants 3 nginx containers running."

---

# 2. Control Plane

The Control Plane remains the brain of the system.

New responsibilities:

- Accept deployment creation requests.
- Store desired state.
- Track current running containers.
- Start the reconciliation process.
- Report deployment status.

The Control Plane should not directly run Docker.

It should communicate with Workers.

---

# 3. Worker

The Worker is responsible for executing container operations.

Responsibilities:

- Receive instructions from the Control Plane.
- Start containers.
- Stop containers.
- Report container status.
- Return success or failure.

Example:

Receive:

```json
{
 "action": "create",
 "image": "nginx"
}
```

↓

Execute:

```bash
docker run nginx
```

↓

Return:

```json
{
 "status": "success"
}
```

---

# 4. Reconciliation Loop

The Reconciliation Loop is the core feature of Phase 2.

Responsibilities:

- Continuously compare desired state and actual state.
- Detect differences.
- Take actions to correct differences.

Example:

Desired:

```
nginx replicas = 3
```

Actual:

```
2 nginx containers running
```

The system should:

```
Create 1 more container
```

until:

```
Desired = Actual
```

---

# 5. State Storage

The system needs to remember desired state.

For Phase 2, you may use:

- In-memory storage
- SQLite database

Example:

```
Deployments:

nginx:
 image: nginx
 replicas: 3
```

The exact storage design is up to you.

---

# Communication Layer

The Control Plane and Worker must communicate over the network.

You may choose:

- REST
- gRPC
- Raw TCP sockets

Be prepared to explain why you chose your approach.

---

# End-to-End Flow

The completed system should work like this:

```
User

↓

mini-kube create deployment nginx --replicas 3

↓

Control Plane receives request

↓

Control Plane stores desired state

↓

Reconciliation Loop checks current state

↓

Worker receives instructions

↓

Docker creates containers

↓

Worker reports current state

↓

Control Plane confirms desired state is achieved
```

---

# Expected Architecture

Example architecture:

```
 User
 |
 |
 v
 mini-kube create deployment
 |
 v
 +----------------+
 | Control Plane |
 +----------------+
 |
 +-----------+-----------+
 | |
 v v
 Desired State Reconciliation Loop
 | |
 +-----------+-----------+
 |
 v
 +----------------+
 | Worker |
 +----------------+
 |
 v
 Docker Containers
 |
 v
 Running Workloads
```

---

# Project Scope

For Phase 2, support:

- Creating deployments.
- Storing desired state.
- Running multiple replicas.
- Tracking running containers.
- Detecting failed containers.
- Automatically recreating missing containers.
- Scaling replica count.

---

# Out of Scope

Do not implement:

- Multiple workers.
- Scheduling.
- Load balancing.
- Service discovery.
- Networking between containers.
- Rolling updates.
- High availability.
- Persistent distributed storage.
- Worker health monitoring.

Those will be implemented in future phases.

---

# Suggested API

Example request:

```
POST /deployments
```

Body:

```json
{
 "name": "nginx",
 "image": "nginx",
 "replicas": 3
}
```

Response:

```json
{
 "status": "created"
}
```

---

Example scaling request:

```
POST /deployments/nginx/scale
```

Body:

```json
{
 "replicas": 5
}
```

Response:

```json
{
 "status": "updated"
}
```

The exact API design is up to you.

---

# Success Criteria

Phase 2 is complete if the following flow works:

```
User

↓

Creates deployment with 3 replicas

↓

Control Plane stores desired state

↓

System creates 3 containers

↓

One container is manually stopped

↓

Reconciliation Loop detects the difference

↓

A replacement container is created

↓

System returns to desired state
```

---

# Deliverables

Each team member should build their own implementation independently.

Each implementation should include:

- Deployment model.
- Control Plane updates.
- Worker updates.
- Desired state storage.
- Reconciliation loop.
- Container tracking.
- Scaling support.

No code sharing until presentations.

---

# Presentation Requirements

Each person should present:

## 1. Architecture

Draw the system.

Explain every component.

---

## 2. Desired State vs Actual State

Explain:

- What is desired state?
- What is actual state?
- How does the system compare them?
- How does the system fix differences?

---

## 3. Request Flow

Walk through exactly what happens after:

```bash
mini-kube create deployment nginx --replicas 3
```

Explain every step until the containers are running.

---

## 4. Design Decisions

Be prepared to answer:

- Why do we need desired state?
- Why use a reconciliation loop?
- How is container state tracked?
- Where is state stored?
- How often does reconciliation run?
- How are failures handled?

---

## 5. Live Demo

Demonstrate:

```bash
mini-kube create deployment nginx --replicas 3
```

Show that:

- The Control Plane receives the request.
- Desired state is stored.
- Containers are created.
- Container failure is detected.
- Replacement containers are created automatically.

---

# Learning Objectives

By completing Phase 2, everyone should understand:

- Desired state systems.
- Reconciliation loops.
- Kubernetes controller concepts.
- Self-healing systems.
- State management.
- Background processing.
- Eventual consistency.
- Container lifecycle management.

---

# Important Rule

The goal is not to produce identical implementations.

Each person should design the system independently.

During the review session, compare:

- Architecture.
- API design.
- State representation.
- Reconciliation logic.
- Error handling.
- Code organization.

Discuss the trade-offs and decide what ideas should be carried forward into the team's shared implementation.

---

# End Goal of Phase 2

Build a system where:

- A user declares the desired number of containers.
- The Control Plane stores that desired state.
- The Worker executes container operations.
- The Reconciliation Loop continuously compares desired and actual state.
- The system automatically fixes failures.

Once this works, we have built the core idea behind Kubernetes controllers and self-healing infrastructure.
