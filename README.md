# `Trellis`

> A calm, concurrent task orchestration runtime built in Go.

`Trellis` is a lightweight runtime for executing, scheduling, and observing background jobs.

It provides a minimal but powerful environment for:

* task execution
* worker pools
* queue management
* job scheduling
* runtime monitoring
* structured observability

`Trellis` focuses on **clear runtime primitives**, **predictable concurrency**, and **operational visibility**.

It is designed as a **minimal orchestration engine**, not a heavy framework.

---

# 🌿 What Is `Trellis`?

`Trellis` is a **task lifecycle runtime**.

It accepts background jobs, queues them, executes them concurrently through workers, and tracks their lifecycle with explicit state transitions.

`Trellis` coordinates:

* task execution
* scheduling
* worker pools
* queue management
* runtime observability

It is **not just a queue** and **not just a worker pool**.

It is the system that **orchestrates both**.

---

# 🧠 Core Concept

At its heart, `Trellis` manages **state and execution**.

Every task moves through a controlled lifecycle:

```
pending → queued → running → completed
                     ↓
                  retrying → failed → dead
```

The runtime ensures transitions are:

* explicit
* observable
* deterministic

This lifecycle model allows `Trellis` to manage execution reliability and system visibility.

---

# 🏗 Architecture Overview

`Trellis` is structured into clear runtime layers.

```
HTTP API
   ↓
Scheduler
   ↓
Queue
   ↓
Worker Pool
   ↓
Task Execution
   ↓
Observability (logs / metrics / traces)
```

A Web UI sits alongside the API to inspect runtime state and task activity.

Each layer has a single responsibility, which keeps the system predictable and maintainable.

---

# ✨ Why `Trellis` Exists

Modern applications depend heavily on background processing:

* sending emails
* generating reports
* refreshing caches
* running ML inference
* performing health checks
* processing long-running tasks

Running this work inside request/response cycles causes:

* slow APIs
* fragile execution
* poor observability

`Trellis` moves this work into a **dedicated orchestration runtime** where tasks are executed safely and observed clearly.

---

# 🚀 Design Principles

`Trellis` follows a few guiding principles.

### Calm Concurrency

Concurrency should be understandable and predictable.

### Explicit State

Task lifecycle transitions must be controlled and validated.

### Minimal Infrastructure

Prefer the Go standard library and avoid unnecessary frameworks.

### Observability First

Systems should be observable while running.

### Clean Boundaries

Each subsystem has a clear responsibility.

`Trellis` favors clarity over cleverness.

---

# ⚙️ Runtime Components

## Task

A **task** represents a unit of background work.

Examples include:

* HTTP health checks
* sending email
* generating reports
* refreshing cached data
* executing ML inference

Each task contains:

* ID
* type
* payload
* status
* timestamps
* optional result
* optional error

Tasks are executed by workers and tracked through their lifecycle.

---

## Queue

The queue stores tasks waiting to be executed.

Responsibilities include:

* enqueue tasks
* dequeue tasks
* support concurrent workers

The initial implementation uses an **in-memory queue**.

Future queue backends may include:

```
Redis
PostgreSQL
distributed queues
```

---

## Worker Pool

Workers execute tasks concurrently.

Workers:

* pull tasks from the queue
* execute task handlers
* update task state
* emit logs and metrics

Workers rely on:

* goroutines
* channels
* context cancellation

Worker concurrency is configurable.

---

## Scheduler

The scheduler creates tasks at defined intervals.

Examples include:

```
every 30 seconds
every minute
cron schedules
```

The scheduler generates tasks and places them into the queue.

Typical scheduled workloads include:

* health checks
* maintenance jobs
* data refresh tasks

---

## Task Handlers

Each task type defines a handler responsible for performing the work.

Conceptually:

```
Execute(ctx context.Context, payload any) error
```

Handlers should:

* respect context cancellation
* emit structured logs
* record metrics

This approach keeps the runtime **generic and extensible**.

---

# 🔍 Observability

Observability is a **first-class feature** of `Trellis`.

Tasks are observable by default, allowing operators to understand what the system is doing while it runs.

---

## Structured Logging

Logging uses:

```
slog
```

Logs include contextual fields such as:

```
task_id
task_type
status
duration
error
```

---

## Metrics

Metrics are exposed through:

```
/metrics
```

Using Prometheus.

Example metrics include:

```
tasks_total
tasks_failed_total
task_duration_seconds
queue_depth
worker_active
```

---

## Tracing

Tracing can be implemented with:

```
OpenTelemetry
```

Tracing allows visibility into:

* task execution
* scheduler events
* queue operations
* HTTP requests

---

# 🌐 HTTP API

`Trellis` exposes a minimal HTTP API for interacting with tasks and runtime state.

Example endpoints:

```
POST   /tasks
GET    /tasks
GET    /tasks/{id}
DELETE /tasks/{id}

GET    /health
GET    /metrics
```

These endpoints allow clients to:

* create tasks
* inspect task state
* monitor runtime health

---

# 🖥 Web UI

`Trellis` includes a Web UI for inspecting the runtime.

The interface provides visibility into:

* active tasks
* task history
* queue depth
* worker activity
* runtime health

The UI is implemented using:

```
GoTH
HTMX
Jaspr components
```

The Web UI interacts with the same API used by external clients.

---

# 🧩 Monitoring as Tasks

Monitoring workloads are implemented as **task types**, not as a separate subsystem.

Example:

```
task type: http_check
```

This task:

* performs an HTTP request
* records latency
* records response status
* emits metrics

This design keeps `Trellis` **generic and extensible**, allowing monitoring to be treated as just another workload.

---

# 📁 Project Structure

```
trellis/

cmd/
    trellis-server/

internal/

    runtime/
        executor.go
        task.go

    queue/
        queue.go

    worker/
        worker.go
        pool.go

    scheduler/
        scheduler.go

    tasks/
        httpcheck/
            task.go
            handler.go

    http/
        router.go
        handlers.go
        middleware.go

    logging/
        logger.go

    observability/
        metrics.go
        tracing.go

web/
    templates/
    components/

pkg/
    response/
```

---

# 🔄 Execution Flow

A typical execution flow in `Trellis` looks like:

```
scheduler
   ↓
enqueue task
   ↓
queue
   ↓
worker
   ↓
task handler
   ↓
logs / metrics / traces
```

This flow ensures tasks are executed asynchronously while remaining observable.

---

# 🧪 Project Status

`Trellis` is currently under active development.

The current runtime includes:

* task lifecycle management
* in-memory task storage
* FIFO task queue
* concurrent worker pool
* handler registry
* HTTP API for task submission and inspection

Future milestones include:

* durable persistence
* task scheduling
* observability instrumentation
* Web UI for runtime inspection

---

# 🎯 Project Goals

`Trellis` aims to:

1. implement a clean orchestration runtime
2. provide strong observability
3. maintain minimal complexity
4. serve as a reference Go backend architecture

The initial runtime is expected to remain small:

```
~1500–2000 LOC
```

---

# 🚫 Non-Goals (for now)

The following are **not part of the initial scope**:

```
distributed cluster mode
alerting systems
workflow DAG engines
multi-node scheduling
```

These may be explored in later versions.

---

# 🌱 Vision

The long-term vision of `Trellis` is to become:

* a minimal yet powerful background job runtime
* a reference architecture for Go concurrency systems
* a clean orchestration layer for modern applications
* a learning platform for building reliable backend systems