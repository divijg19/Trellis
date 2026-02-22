# `Trellis`

> A calm, concurrent task runtime built in Go.

`Trellis` is a background task orchestration engine.
It accepts tasks, schedules them, executes them concurrently through workers, retries failures, and tracks the full lifecycle of each task with clarity and reliability.

---

## 🌿 What Is `Trellis`?

`Trellis` is a concurrent task lifecycle manager.

It is designed to:

* Accept background jobs via API
* Queue and schedule work
* Execute tasks concurrently
* Retry failures with policy
* Persist state across restarts
* Provide observability into task execution

`Trellis` is not just a queue.
It is not just a worker pool.
It is the system that coordinates both.

---

## 🧠 Core Concept

At its heart, `Trellis` manages **state and execution**.

Every task moves through a lifecycle:

```
pending → queued → running → completed
                     ↓
                  retrying → failed → dead
```

`Trellis` ensures transitions are controlled, observable, and durable.

---

## 🏗 Architecture Overview

`Trellis` is structured in clear layers:

```
domain      → Task definitions and lifecycle rules
runtime     → Orchestration engine
queue       → Task transport
worker      → Concurrent execution
storage     → Persistence layer
api         → HTTP interface
web         → UI (GoTH + Jaspr)
```

This separation keeps `Trellis` predictable and maintainable.

---

## ✨ Why `Trellis` Exists

Modern applications require reliable background processing:

* Email sending
* Report generation
* Image processing
* External API calls
* AI pipelines
* Long-running computations

Blocking request/response cycles is inefficient and fragile.

`Trellis` moves that work into a structured background runtime.

---

## 🚀 Design Principles

`Trellis` follows a few core principles:

* Calm concurrency
* Explicit state transitions
* Minimal magic
* Clear failure handling
* Graceful shutdown
* Observability first
* Clean architecture boundaries

`Trellis` favors clarity over cleverness.

---

## 🔄 Task Lifecycle

Each task in `Trellis` contains:

* Identifier
* Type
* Payload
* Status
* Retry metadata
* Timestamps
* Error information
* Optional result

State transitions are explicit and validated.

Retries follow a defined policy (e.g., exponential backoff).

---

## ⚙️ Concurrency Model

`Trellis` uses:

* Goroutines for workers
* Channels for queue transport
* Context propagation for cancellation
* Configurable worker concurrency

Workers execute handlers registered by task type.

---

## 🗄 Persistence

`Trellis` persists tasks to durable storage.

On restart, `Trellis`:

* Recovers unfinished tasks
* Requeues eligible work
* Maintains consistency

Durability is a core feature, not an afterthought.

---

## 📊 Observability

`Trellis` exposes:

* Task states
* Retry counts
* Worker status
* Metrics endpoint
* Structured logs

Systems should be understandable while running.

---

## 🛠 CLI (Initial)

```
trellis server
trellis worker --concurrency 5
```

`Trellis` is designed to support single-process and distributed worker modes.

---

## 🧪 Project Status

`Trellis` is under active development.

Initial focus:

1. Domain model
2. In-memory runtime
3. Worker pool
4. HTTP API
5. Persistence
6. Scheduler
7. Observability
8. UI (GoTH + Jaspr)

---

## 🎯 Long-Term Vision

`Trellis` aims to be:

* A minimal yet powerful background runtime
* A learning platform for concurrency mastery
* A production-ready orchestration core
* A foundational system within a broader ecosystem