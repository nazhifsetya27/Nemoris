# **Week 1 — Core Stabilization**

# Sprint Goal

Stabilize current MVP foundation before adding production complexity.

Current priority:

> Prevent architecture debt while system is still small.

---

# EPIC 1 — Architecture Safety Refactor

## Tasks

- Move database migration out of `main.go`
- Create `internal/database/migrate.go`
- Keep `main.go` only for boot wiring
- Ensure migration callable from one entry point only

## Done Criteria

- `main.go` no longer contains migration logic
- startup sequence becomes:

```
config → database → migrate → routes → server
```

## Priority

Critical

---

# EPIC 2 — Configuration Centralization

## Tasks

- Create `internal/config/config.go`
- Move all env loading there
- Remove direct `os.Getenv()` from service layer
- Expose config struct globally safe

## Done Criteria

No service package directly reads environment variables

## Priority

Critical

---

# EPIC 3 — WhatsApp Payload Normalization Layer

## Tasks

- Create `internal/whatsapp/payload.go`
- Normalize WAHA webhook payload
- Separate transport payload from internal message model
- Support future payload version changes safely

## Done Criteria

Handler never directly trusts WAHA raw payload

## Priority

Critical

---

# EPIC 4 — Safe Outbound Messaging Foundation

## Tasks

- Create `internal/whatsapp/send.go`
- Add outbound sender wrapper
- Block empty outbound send
- Block self-message send
- Log WAHA response

## Done Criteria

All outbound messages pass through single sender layer

## Priority

Critical

---

# EPIC 5 — Structured Logging

## Tasks

- Create `internal/utils/logger.go`
- Replace raw `log.Println`
- Add event categories:
  - inbound
  - outbound
  - db
  - scheduler
  - ai

## Done Criteria

Logs readable per request flow

## Priority

High

---

# EPIC 6 — Database Performance Protection

## Tasks

- Add index for duplicate check
- Add index for reminder time lookup
- Validate GORM migration still safe

## Done Criteria

Indexes exist:

```
messages(from, body, created_at)
reminders(remind_at)
```

## Priority

High

---

# Sprint Exit Criteria

Sprint complete only if:

- startup clean
- architecture safe
- no direct env leak
- outbound sender exists

---

# Sprint Risk

If skipped:

future AI + scheduler becomes unstable.

# **Week 2 — Reminder Engine Production**

# Sprint Goal

Turn saved reminders into actual delivery system.

---

# EPIC 1 — Scheduler Foundation

## Tasks

- Create `internal/scheduler`
- Add polling loop
- Execute every minute
- Load due reminders safely

## Done Criteria

Scheduler runs independently

## Priority

Critical

---

# EPIC 2 — Reminder Status Lifecycle

## Tasks

Add fields:

- status
- sent_at
- retry_count

States:

- pending
- sent
- failed
- retrying

## Done Criteria

Reminder lifecycle visible in DB

## Priority

Critical

---

# EPIC 3 — Reminder Worker Flow

## Tasks

Build flow:

```
scheduler → service → repository → whatsapp sender
```

## Done Criteria

Reminder sent through service layer only

## Priority

Critical

---

# EPIC 4 — Retry Logic

## Tasks

- max retry 3
- retry delay
- fail final safely

## Done Criteria

Failed WA send does not disappear silently

## Priority

High

---

# EPIC 5 — Duplicate Delivery Protection

## Tasks

- atomic status update before send
- lock reminder row

## Done Criteria

Same reminder cannot send twice

## Priority

Critical

---

# EPIC 6 — Reminder Read API

## Tasks

- add `/reminders`
- filter by sender

## Done Criteria

Users can inspect reminder state

## Priority

Medium

---

# Sprint Exit Criteria

Reminder actually reaches WhatsApp safely.

---

# Sprint Risk

Without this sprint:

Nemoris only stores memory but cannot act.

# Week 3 — AI + Multilingual

# Sprint Goal

Transform parser from regex bot into bilingual assistant.

---

# EPIC 1 — AI Layer Separation

## Tasks

Create:

```
internal/ai/
```

Files:

- intent.go
- parser.go
- language.go

## Done Criteria

AI logic leaves service package

## Priority

Critical

---

# EPIC 2 — Language Detection

## Tasks

Support:

- Indonesian
- English

Examples:

- remind me tomorrow
- ingatkan saya besok

## Done Criteria

Incoming message classified before parsing

## Priority

Critical

---

# EPIC 3 — Dual Parsing Strategy

## Tasks

Order:

1 rule parser

2 AI fallback

## Done Criteria

Cheap messages avoid AI call

## Priority

Critical

---

# EPIC 4 — Intent Engine

## Tasks

Support intents:

- create_reminder
- list_reminders
- store_memory
- ignore_smalltalk

## Done Criteria

Intent output standardized

## Priority

Critical

---

# EPIC 5 — JSON AI Contract

## Tasks

Force AI output:

```
{
  "intent":"create_reminder",
  "task":"pay electricity",
  "time":"2026-03-08T20:00:00",
  "lang":"id"
}
```

## Done Criteria

No free text AI output allowed

## Priority

Critical

---

# EPIC 6 — Clarification Logic

## Tasks

Ask clarification when confidence low

## Done Criteria

Ambiguous message never silently misparsed

## Priority

High

---

# Sprint Exit Criteria

System understands two languages safely.

---

# Sprint Risk

Without confidence control:

AI causes wrong reminders.

# Week 4 — Production Hardening

# Sprint Goal

Make Nemoris survivable in real user traffic.

---

# EPIC 1 — Webhook Security

## Tasks

- add webhook secret validation
- reject unknown source

## Done Criteria

Unauthorized webhook blocked

## Priority

Critical

---

# EPIC 2 — Rate Limiting

## Tasks

- sender request limit
- burst protection

## Done Criteria

Spam cannot overload backend

## Priority

Critical

---

# EPIC 3 — Panic Recovery

## Tasks

- middleware recover panic
- log panic safely

## Done Criteria

Crash does not kill process

## Priority

Critical

---

# EPIC 4 — Failed Job Tracking

## Tasks

Create failed_jobs tracking

## Done Criteria

Failures visible for investigation

## Priority

High

---

# EPIC 5 — Health Check Upgrade

## Tasks

Health endpoint checks:

- DB
- Redis
- WAHA

## Done Criteria

Health reflects real service condition

## Priority

Critical

---

# EPIC 6 — WhatsApp Compliance Layer

## Tasks

- delay outbound
- anti-loop protection
- safe send pacing

## Done Criteria

Ban risk reduced

## Priority

Critical

---

# EPIC 7 — Backup Strategy

## Tasks

- pg_dump daily
- restore test

## Done Criteria

Recovery proven

## Priority

High

---

# EPIC 8 — Production Compose Split

## Tasks

Separate:

- dev compose
- prod compose

## Done Criteria

Deployment predictable

## Priority

Medium

---

# Sprint Exit Criteria

System survives production stress.

---

# Sprint Risk

Without this sprint:

system works but breaks under real users.
