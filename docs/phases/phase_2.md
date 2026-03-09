# **Phase 2 — Reliability + Product Intelligence**

> Goal:
>
> Transform NEMORIS from stable MVP into production-shaped core system:
>
> stronger scheduler, cleaner multilingual replies, first memory layer, and observability.

This phase assumes **Phase 1 (Sprint 1–4) is complete**:

- webhook stable
- reminders stable
- bilingual parser stable
- retry lifecycle stable
- production layering enforced

References from current architecture already support this direction: scheduler remains in-process, memory intent exists but not implemented, localized reply builder is planned for future expansion
architecture
project_instructions

---

# Sprint 5 — Scheduler Hardening

## Sprint Goal

Make reminder execution reliable under restart, timing drift, and future multi-worker growth.

---

## EPIC 1 — Scheduler Metrics

### Tasks

Create:

```
internal/scheduler/metrics.go
```

Track:

- total checked reminders
- due reminders
- sent reminders
- retry reminders
- failed reminders
- execution duration

### Done Criteria

Scheduler prints measurable execution stats every cycle.

### Priority

Critical

---

## EPIC 2 — Claim Safety Upgrade

### Tasks

Strengthen reminder claiming query.

Use DB-safe lock:

```
FORUPDATE SKIP LOCKED
```

inside repository claim flow.

### Done Criteria

Same reminder cannot be claimed twice even under concurrent worker.

### Priority

Critical

---

## EPIC 3 — Scheduler Drift Detection

### Tasks

Detect late execution:

- compare current time vs expected interval
- log drift warning if delayed

### Done Criteria

Long pauses visible in logs.

### Priority

High

---

## EPIC 4 — Scheduler Recovery Logging

### Tasks

On startup:

- count pending reminders
- count overdue reminders

### Done Criteria

Restart immediately shows backlog status.

### Priority

High

---

## EPIC 5 — Reminder Execution Profiling

### Tasks

Measure:

- claim duration
- send duration
- update duration

### Done Criteria

Slow path visible before scaling.

### Priority

Medium

---

## Sprint Exit Criteria

Scheduler survives:

- restart
- drift
- duplicate execution risk

---

## Sprint Risk

Without this sprint:

Reminder reliability degrades silently as usage grows.

---

# Sprint 6 — Localized Reply System

## Sprint Goal

Separate reply language from business logic cleanly.

---

## EPIC 1 — Reply Builder Layer

### Tasks

Create:

```
internal/i18n/
 ├── reply.go
 ├── id.go
 ├── en.go
```

### Done Criteria

All user replies generated through one layer.

### Priority

Critical

---

## EPIC 2 — Canonical Intent → Localized Reply

### Tasks

Map canonical intent:

```
create_reminder
list_reminders
store_memory
unknown
```

to localized responses.

### Done Criteria

Service no longer returns hardcoded text.

### Priority

Critical

---

## EPIC 3 — Reply Templates

### Tasks

Support placeholders:

```
task
time
count
```

Example:

```
Okay, I'll remind you tomorrow at 8 PM.
```

```
Siap, saya ingatkan besok jam 8 malam.
```

### Done Criteria

Dynamic reply text centralized.

### Priority

High

---

## EPIC 4 — Language Fallback Safety

### Tasks

If language unknown:

fallback:

```
en
```

### Done Criteria

No empty reply caused by detection miss.

### Priority

High

---

## EPIC 5 — Future Language Expansion Hook

### Tasks

Prepare:

```
RegisterLanguage()
```

pattern for future languages.

### Done Criteria

Adding language does not require service rewrite.

### Priority

Medium

---

## Sprint Exit Criteria

Replies fully detached from parser/business logic.

---

## Sprint Risk

Without this sprint:

Multilingual grows messy fast.

---

# Sprint 7 — Memory Domain MVP

## Sprint Goal

Start NEMORIS core identity beyond reminders.

---

## EPIC 1 — Memory Table

### Tasks

Create:

```
memories
```

Fields:

```
id UUID
from TEXT
content TEXT
lang TEXT
created_atTIMESTAMP
```

### Done Criteria

Memory persistence available.

### Priority

Critical

---

## EPIC 2 — Memory Model + Repository

### Tasks

Add:

```
internal/model/memory.go
internal/repository/memory_repository.go
```

### Done Criteria

Repository fully isolated.

### Priority

Critical

---

## EPIC 3 — Store Memory Intent Activation

### Tasks

Activate existing canonical intent:

```
store_memory
```

### Done Criteria

Messages classified as memory get saved.

### Priority

Critical

---

## EPIC 4 — Basic Memory Confirmation Reply

### Tasks

Reply example:

```
Noted, I'll remember that.
```

```
Baik, saya simpan.
```

### Done Criteria

User sees memory saved.

### Priority

High

---

## EPIC 5 — Memory Read Endpoint

### Tasks

Add:

```
/memories?from=
```

### Done Criteria

Stored memory inspectable.

### Priority

Medium

---

## Sprint Exit Criteria

NEMORIS stores non-reminder information safely.

---

## Sprint Risk

Without this sprint:

Product remains reminder-only.

---

# Sprint 8 — Observability + Failure Intelligence

## Sprint Goal

Make production behavior diagnosable.

---

## EPIC 1 — Request Correlation ID

### Tasks

Generate request id per inbound webhook.

Propagate through:

- handler
- service
- repository
- outbound

### Done Criteria

Single request traceable end-to-end.

### Priority

Critical

---

## EPIC 2 — Structured Failure Types

### Tasks

Add failure type field:

```
failure_type
```

Examples:

- waha_timeout
- invalid_target
- db_lock_fail
- unknown_send_error

### Done Criteria

Failure reason visible.

### Priority

Critical

---

## EPIC 3 — Failed Reminder Analytics

### Tasks

Expand failed reminder API.

Add filters:

```
status
failure_type
date
```

### Done Criteria

Failure patterns inspectable.

### Priority

High

---

## EPIC 4 — Health Check Expansion

### Tasks

Health endpoint returns:

- DB latency
- WAHA latency
- scheduler last tick

### Done Criteria

Health reflects live operational state.

### Priority

Critical

---

## EPIC 5 — Log Noise Cleanup

### Tasks

Separate:

- info
- warning
- error

### Done Criteria

Logs readable under real traffic.

### Priority

High

---

## Sprint Exit Criteria

Production issues diagnosable without guessing.

---

## Sprint Risk

Without this sprint:

Real failures become invisible.

---

# Phase 2 Exit Criteria

Phase 2 completes only if:

- scheduler hardened
- replies localized cleanly
- memory saved
- failures diagnosable

---

# After Phase 2 → Phase 3 Ready

Then safe to begin:

- recurring reminders
- voice notes
- semantic memory retrieval
- Redis queue evolution
- external AI fallback

---

# Strong Principle for Phase 2

> Do not add external intelligence before internal reliability becomes boring.

Because:

> boring systems survive 🚀
