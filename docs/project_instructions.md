# NEMORIS Project Description

*Last updated: 2026-03-08*

NEMORIS is a production-oriented WhatsApp-native AI memory assistant built with:

- Go backend
- WAHA Free (WhatsApp HTTP API)
- PostgreSQL
- Redis
- Docker Compose
- GORM (Go ORM for database persistence)

Primary goal:

Build a stable, production-safe system that receives WhatsApp messages, interprets user intent, stores reminders/tasks/memory, and sends scheduled reminders safely.

## Multilingual Support (MVP)

NEMORIS supports **2 languages from MVP stage**: Indonesian (`id`) and English (`en`).

- **Bilingual user input** — Indonesian and English natural language.
- **Unified internal intent processing** — NEMORIS accepts both Indonesian and English user messages while converting them into one language-neutral internal command model.
- **Localized Reply** — Responses generated via Localized Reply Builder based on detected language.

---

# Development Environment Baseline

NEMORIS development currently runs on lightweight Apple Silicon hardware.

## Local Machine

- Machine: MacBook Air M1
- RAM: 8GB
- Storage: 256GB SSD
- Architecture: Apple Silicon (ARM64)

## Container Runtime

Using Colima instead of Docker Desktop:

```bash
colima start --cpu 2 --memory 2 --vz-rosetta --dns 8.8.8.8
```

## Why This Configuration Exists

- `--cpu 2` keeps macOS responsive during development
- `--memory 2` prevents Docker VM from exhausting 8GB RAM
- `--vz-rosetta` improves x86 image compatibility on Apple Silicon
- `--dns 8.8.8.8` avoids Docker DNS instability on macOS

## Docker Versions

```bash
docker -v
Docker version 29.2.1
```

```bash
docker compose -v
Docker version 29.2.1
```

## Current Local Containers

```bash
docker ps
```

```text
WAHA          → port 3000
PostgreSQL 15 → port 5432
Redis 7       → port 6379
```

## Current Running Containers

```text
waha
nemoris-postgres-1
nemoris-redis-1
```

## Development Rule

During MVP:

- infrastructure runs in containers
- Go backend runs locally on host machine first

Reason:

- faster iteration
- easier debugging
- simpler logs
- lower RAM usage

## Local Resource Safety Rule

On MacBook Air 8GB RAM:

Never overload system with:

- heavy Docker stacks
- large AI models
- unnecessary parallel services

Safe development baseline:

```text
WAHA + PostgreSQL + Redis + Go backend only
```

Stop unused containers immediately.

## Infrastructure Priority Order Under Memory Pressure

Keep alive first:

1. PostgreSQL
2. Redis
3. WAHA

Restartable safely:

4. Go backend

---

## Core Engineering Principles

- Prioritize production safety over speed

- Avoid anything that risks WhatsApp number bans

- Never assume webhook payloads before inspecting actual payloads

- Prefer incremental development:

  1. receive
  2. log
  3. validate
  4. persist
  5. send

- Every feature must be testable locally before touching WAHA live messaging

---

## Current Technical Stack

- Backend language: Go (Go 1.23+)
- API style: standard net/http first, then evolve toward Gin only when needed
- Messaging gateway: WAHA Free Docker
- Database: PostgreSQL 15
- ORM: GORM
- Queue / scheduler: Redis 7
- Container runtime: Docker Compose + Colima on Apple Silicon

---

## Database Layer Rules (GORM + PostgreSQL)

- Use GORM as primary ORM during MVP for faster safe iteration.

- PostgreSQL remains source of truth.

- GORM usage is allowed only inside repository layer.

Never allow:

```text
handler → gorm direct call
service → gorm direct call
```

Always:

```text
handler → service → repository → gorm → postgres
```

- `database/` owns single shared GORM connection lifecycle.

- `repository/` owns all `Create`, `Find`, `Where`, `Updates`, `Delete` operations.

- `service/` must never know GORM internals.

- Use `AutoMigrate()` only during MVP local development.

- When schema stabilizes, move to explicit migration system.

## GORM Safety Rules

- Define one model struct per domain entity.

- Separate transport structs from persistence structs.

Correct example:

```text
IncomingMessage ≠ Message
```

Because:

- incoming webhook payload belongs to transport layer
- DB entity belongs to persistence layer

## Model Separation Rule

Use split files inside model:

```text
internal/model/
 ├── incoming_message.go
 ├── message.go
 ├── reminder.go
```

Never mix many unrelated structs in one large file once domain grows.

## Database Connection Rule

Only one connection source:

```text
internal/database/postgres.go
```

Never create DB connection elsewhere.

## Environment Safety Rule

During MVP:

connection string may be simple.

After first stable persistence:

move credentials into:

```text
.env
internal/config/
```

Never keep production secrets hardcoded.

---

## Code Philosophy

Guide explanations assuming developer is strong in JavaScript but beginner in Go.

Always explain:

- Go concepts using JavaScript comparisons
- folder structure reasoning
- why a pattern is production-safe
- common beginner mistakes

Avoid overengineering early.

When explaining ORM:

Always compare:

```text
Sequelize ↔ GORM
```

So migration from JavaScript mental model stays easy.

---

## Backend Architecture (Current)

Standard Go production layout:

```text
backend/
 ├── cmd/
 │   └── api/
 │       └── main.go
 │
 ├── internal/
 │   ├── ai/           # intent extraction, language detection, reminder parsing (rule-based)
 │   ├── config/       # environment loading, app config, secrets mapping
 │   ├── database/     # postgres connection, migrations, db bootstrap
 │   ├── handler/      # HTTP handlers / webhook entrypoints
 │   ├── middleware/   # recover, rate limit, webhook auth
 │   ├── model/        # shared structs and domain models
 │   ├── repository/   # database access layer via GORM
 │   ├── scheduler/    # in-process 60s ticker → ProcessDueReminders
 │   ├── service/      # orchestration between modules
 │   ├── utils/        # shared helpers, logging
 │   └── whatsapp/     # WAHA integration, send/receive, LID resolution
 │
 ├── Dockerfile
 └── go.mod
```

### Multilingual Components

- **ai/language.go** — Detects `id` or `en` from incoming message; runs before Intent Parsing.
- **ai/intent.go** — Canonical intent detection.
- **Localized reply builder** — Planned for Phase 2; currently uses simple strings.

---

## Rules

- `cmd/api/main.go` stays minimal
  only boot server, wire dependencies, register routes.

- `handler/` only:
  receive request → validate → call service → return response.

- `service/` contains business orchestration:
  reminder creation, memory processing, AI flow coordination.

- `repository/` handles all persistence:
  GORM queries only here.

- `whatsapp/` must isolate WAHA-specific logic:
  webhook payload parsing, outbound send protection, retry safety.

- `scheduler/` owns all delayed execution:
  reminders, recurring tasks, future notifications.

- `ai/` must remain isolated from reminder logic:
  AI returns structured intent only (language, intent, task, time); Canonical output is language-independent; business logic stays elsewhere.

- `utils/` only for generic reusable helpers:
  avoid putting business logic here.

- `model/` contains domain-level structs shared across layers.

- `database/` owns connection lifecycle:
  app should never create raw DB connections elsewhere.

---

## Production Safety Rule

Never allow:

```text
handler → repository direct call
```

Always:

```text
handler → service → repository
```

This keeps logic maintainable when Nemoris grows.

---

## WAHA Safety Rules

Before any outbound send:

- inspect full webhook payload first
- ignore empty payloads
- ignore self messages
- ignore duplicate events
- never auto-reply broadly
- reply only to explicit safe commands during development

Never create loops.

---

## Instruction Maintenance Rule

Whenever a new stack, library, infrastructure component, framework, ORM, queue system, AI provider, deployment pattern, or architectural dependency is adopted:

ChatGPT must explicitly remind:

```text
Update project instructions now.
```

Because architecture memory must stay synchronized with actual stack.

Never continue long implementation after major stack adoption without first suggesting instruction update.

Examples:

- adding GORM
- adding Gin
- adding Asynq
- adding OpenAI SDK
- adding migration tool
- adding worker container

Instruction updates must happen immediately after stack decision.

---

## Chat Separation Rules

Each new chat should have one clear purpose only.

Suggested chat naming:

- NEMORIS / Setup
- NEMORIS / Debug Backend
- NEMORIS / Database Design
- NEMORIS / WAHA Integration
- NEMORIS / Scheduler
- NEMORIS / AI Parsing
- NEMORIS / Production Deployment

Each chat must stay focused on one scope.

---

## Response Style Required

Responses should be:

- beginner-friendly
- step-by-step
- explicit
- production-minded
- honest about risks
- explain why before how

Always prefer safe incremental implementation.

---

## Documentation Standard

When generating documentation:

- use markdown
- use mermaid diagrams when architecture is discussed
- use production-grade naming
- separate MVP vs future phases clearly

## Terminology (Multilingual)

Use consistent terminology everywhere:

- **Language Detection** — detecting `id` or `en` from incoming message
- **Canonical Intent** — language-independent internal intent (e.g. `create_reminder`, never `buat_pengingat`)
- **Localized Reply** — response text generated in user's language

---

## Long-Term Goal

NEMORIS should become:

> A reliable invisible WhatsApp memory layer where users trust that one message is enough and the system never forgets.

---

## Very Important Local Constraint

Because development happens on **8GB RAM machine**, every architectural choice must stay lightweight first.

Meaning:

- prefer net/http before framework
- prefer simple GORM before complex abstraction
- avoid early distributed patterns
- avoid premature worker explosion

Production sophistication comes later, after stable MVP.
