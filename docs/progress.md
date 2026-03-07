# NEMORIS Progress Summary

## Current Build Status

This summary reflects what has already been completed in the current MVP backend build of **NEMORIS**.

NEMORIS supports **2 languages from MVP stage**: Indonesian (`id`) and English (`en`). Users send messages in either language; the system converts them into a language-neutral Canonical Intent and returns Localized Replies.

---

# Phase 0 — Infrastructure Baseline ✅

## Local runtime established

- Docker Compose running successfully
- WAHA container active on port `3000`
- PostgreSQL container active on port `5432`
- Redis container active on port `6379`
- Go backend runs locally outside Docker (safe MVP choice for M1 Mac)

## Infrastructure safety decisions

- WAHA pinned with DNS override
- Restart policy enabled for WAHA
- PostgreSQL and Redis isolated in containers
- Minimal safe container footprint maintained

## Current compose stack

```text
WAHA + PostgreSQL + Redis + Local Go Backend
```

---

# Phase 1 — Backend Core Bootstrapped ✅

## Go project initialized

- `go.mod` created
- Go version `1.23`
- Dockerfile prepared for future backend containerization

## Folder structure established

```text
backend/
├── cmd/api/main.go
├── internal/
│   ├── database/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   ├── service/
│   └── utils (planned)
```

## Dependency baseline added

- GORM
- PostgreSQL driver
- dotenv loader

---

# Phase 2 — Database Layer Working ✅

## PostgreSQL connection completed

- Shared connection centralized in:

```text
internal/database/postgres.go
```

## GORM active

- AutoMigrate enabled for MVP
- Connected successfully using env-based DSN

## Active tables

- `messages`
- `reminders`

## Safety rule applied

- GORM usage restricted to repository layer only

---

# Phase 3 — Message Intake Pipeline Working ✅

## HTTP server running

Active routes:

- `/health`
- `/webhook`

## Health endpoint verified

Returns:

```json
{ "status": "ok" }
```

## Webhook intake active

Current flow:

```text
WAHA → handler → service → repository → postgres
```

## Incoming payload model split correctly

- `IncomingMessage`
- `Message`

This avoids transport/persistence coupling.

---

# Phase 4 — Message Persistence Working ✅

## Incoming messages saved

Each valid webhook message is persisted into PostgreSQL.

Stored fields:

- sender
- body
- created_at

## Duplicate protection added

Current rule:

```text
same sender + same body within 10 seconds = ignored
```

## Self-message protection added

Current rule:

```text
BOT_NUMBER ignored
```

This prevents future WhatsApp loops.

---

# Phase 5 — Rule-Based Reminder Parsing Working ✅

## First parser implemented

Supported commands (bilingual):

**English:**

```text
remind me to pay electricity tomorrow 8pm
```

**Indonesian:**

```text
ingatkan saya untuk bayar listrik besok jam 8 malam
```

Both map internally to the same Canonical Intent.

## Parser extracts

- task
- relative time text

## Reminder persistence active

Saved into reminders table. `language` field stored for analytics and localized delivery; business logic remains language-neutral.

---

# Phase 6 — Reminder Listing Working ✅

## `/reminders` endpoint working

Verified output example:

```json
[
  {
    "ID": 1,
    "From": "628123",
    "Task": "pay electricity",
    "RawTime": "tomorrow 8pm",
    "RemindAt": "2026-03-08T20:00:00+07:00"
  }
]
```

## Current achievement

Reminder now reaches structured timestamp storage.

This is the first true persistence milestone.

---

# Multilingual Foundation (MVP) ✅

- **Language Detection** — Detects `id` or `en` from incoming message before parsing.
- **Canonical Intent** — Internal model is language-independent (e.g. `create_reminder`, never `buat_pengingat`).
- **Localized Reply** — Responses generated in user's language (e.g. "Reminder noted" / "Pengingat disimpan").

# Current MVP Capability Right Now ✅

A WhatsApp-style message can already:

```text
enter system → language detection → parse → canonical intent → save → localized reply → list
```

Meaning core memory pipeline is alive.

---

# Immediate Next Development Phases 🚀

---

# Phase 7 — Scheduler Engine (Next Priority)

## Goal

Trigger reminders automatically when due.

## Build target

```text
scheduler checks reminders every minute
```

## First safe MVP

- polling scheduler
- no Redis yet
- local goroutine scheduler

## Output target

```text
Reminder due: pay electricity
```

---

# Phase 8 — WAHA Outbound Messaging

## Goal

Send actual reminder back to WhatsApp.

## Required build

- WAHA sendText integration
- safe outbound payload validation
- API key protected requests

## Critical safety

Never send automatically until scheduler fully trusted.

---

# Phase 9 — Strong Time Parsing

## Replace current parser

Current parser is rule-only.

Next target:

Support:

```text
next monday 7am
in 2 hours
15 march 9pm
```

## Recommended package

Use deterministic parser before LLM.

---

# Phase 10 — Recurring Reminders

## Target

Support:

```text
every monday
monthly
weekly
```

## Required schema extension

Add:

- recurrence
- status
- next_run_at

---

# Phase 11 — AI Intent Layer

## Goal

Move beyond rule parsing.

## AI handles

- reminder intent
- memory intent
- query intent
- note storage

## Example (bilingual)

**English:**

```text
remember my passport expires in december
```

**Indonesian:**

```text
ingatkan saya paspor saya kadaluarsa bulan desember
```

## Output contract

AI must return: `language`, `intent`, `task`, `time` — Canonical output independent of source language.

---

# Phase 12 — Memory System

## Separate from reminders

New table:

```text
memories
```

Supports:

- facts
- long-term notes
- retrieval

---

# Phase 13 — Redis Delayed Jobs

## Upgrade scheduler

Replace polling with queue.

Recommended:

```text
Asynq
```

Only after local scheduler proven stable.

---

# Phase 14 — Production Hardening

## Add

- structured logging
- env config package
- migration tool
- backup policy
- WAHA webhook auth

---

# Long-Term Product Phases 🌱

## Phase A

- voice note transcription
- media understanding

## Phase B

- calendar sync
- proactive reminders

## Phase C

- invisible AI memory companion

## Future Multilingual Expansion

- Regional language support
- Voice multilingual parsing

---

# Core Principle Preserved

> one message sent, system never forgets
