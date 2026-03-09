# **Phase 3 — Scalable Intelligence + Durable Execution**

> Goal:
>
> Move NEMORIS from single-node smart assistant into scalable intelligent memory infrastructure.

This phase should begin only after:

- scheduler hardened
- reply builder stable
- memory domain active
- observability working

Phase 3 matches the future expansion already implied in your architecture docs: Redis queue, recurring tasks, memory retrieval, and deeper intelligence layer
architecture

---

# Sprint 9 — Recurring Reminder Engine

## Sprint Goal

Upgrade reminders from one-time execution into repeatable time logic.

---

## EPIC 1 — Recurrence Schema

### Tasks

Add fields to reminders:

```
recurrence_type
recurrence_interval
next_run_at
```

Supported MVP values:

```
daily
weekly
monthly
```

### Done Criteria

Reminder schema supports recurring lifecycle.

### Priority

Critical

---

## EPIC 2 — Recurrence Calculator

### Tasks

Create:

```
internal/service/recurrence_service.go
```

Responsibilities:

- calculate next occurrence
- preserve timezone
- prevent past scheduling

### Done Criteria

Next run generated after successful send.

### Priority

Critical

---

## EPIC 3 — Scheduler Recurring Flow

### Tasks

After send success:

```
sent → generate next pending
```

### Done Criteria

Recurring reminder never disappears after first execution.

### Priority

Critical

---

## EPIC 4 — Natural Language Recurrence Parsing

### Tasks

Support:

English:

```
every monday
every day at 8
```

Indonesian:

```
setiap senin
tiap hari jam 8
```

### Done Criteria

Parser emits recurrence safely.

### Priority

High

---

## EPIC 5 — Recurrence Safety Guard

### Tasks

Reject ambiguous recurrence:

```
sometimes every month maybe
```

### Done Criteria

Unsafe recurrence asks clarification.

### Priority

High

---

## Sprint Exit Criteria

Recurring reminders survive repeated cycles safely.

---

## Sprint Risk

Without guardrails:

Recurring reminders create silent infinite mistakes.

---

# Sprint 10 — Redis Queue Evolution

## Sprint Goal

Move reminder execution toward durable queue architecture.

---

## EPIC 1 — Redis Responsibility Activation

### Tasks

Redis becomes active for:

```
execution lock
delayed dispatch marker
```

Not full worker migration yet.

### Done Criteria

Redis now has operational responsibility.

### Priority

Critical

---

## EPIC 2 — Queue Adapter Layer

### Tasks

Create:

```
internal/queue/
```

Files:

```
enqueue.go
dequeue.go
lock.go
```

### Done Criteria

Queue isolated from scheduler logic.

### Priority

Critical

---

## EPIC 3 — Scheduler Hybrid Mode

### Tasks

Flow:

```
scheduler detects due
→ enqueue
→ worker executes
```

### Done Criteria

Execution separated from detection.

### Priority

Critical

---

## EPIC 4 — Retry via Queue

### Tasks

Retry delay controlled through queue.

### Done Criteria

Retry timing independent from main loop.

### Priority

High

---

## EPIC 5 — Queue Failure Visibility

### Tasks

Track:

- enqueue fail
- dequeue fail
- lock timeout

### Done Criteria

Queue failures visible.

### Priority

High

---

## Sprint Exit Criteria

Reminder execution survives backend timing fluctuation.

---

## Sprint Risk

Without queue separation:

Scaling causes scheduler stress.

---

# Sprint 11 — Memory Retrieval Engine

## Sprint Goal

Turn stored memory into usable assistant capability.

---

## EPIC 1 — Memory Query Intent

### Tasks

Add intent:

```
retrieve_memory
```

### Done Criteria

System distinguishes ask vs save.

### Priority

Critical

---

## EPIC 2 — Retrieval Repository

### Tasks

Add search:

- latest memory
- keyword match
- sender scoped

### Done Criteria

Basic retrieval works safely.

### Priority

Critical

---

## EPIC 3 — Bilingual Query Support

### Tasks

Examples:

English:

```
what is my bank account?
```

Indonesian:

```
berapa nomor rekening saya?
```

### Done Criteria

Canonical retrieval intent stable.

### Priority

Critical

---

## EPIC 4 — Retrieval Reply Builder

### Tasks

Reply in detected language.

### Done Criteria

Retrieved memory localized.

### Priority

High

---

## EPIC 5 — Safe Retrieval Guard

### Tasks

If multiple match:

ask clarification.

### Done Criteria

Wrong memory not returned silently.

### Priority

High

---

## Sprint Exit Criteria

NEMORIS can remember and answer.

---

## Sprint Risk

Without safe retrieval:

memory becomes unreliable.

---

# Sprint 12 — AI Gateway Foundation

## Sprint Goal

Prepare controlled external intelligence without losing architecture safety.

---

## EPIC 1 — AI Provider Isolation

### Tasks

Create:

```
internal/ai/provider/
```

### Done Criteria

External AI isolated from parser core.

### Priority

Critical

---

## EPIC 2 — Rule First, AI Fallback

### Tasks

Flow:

```
rule parser
→ if low confidence
→ external AI fallback
```

### Done Criteria

Cheap path remains default.

### Priority

Critical

---

## EPIC 3 — JSON Contract Enforcement

### Tasks

AI output must always return:

```
{
  "intent":"",
  "task":"",
  "time":"",
  "lang":""
}
```

### Done Criteria

No raw prose allowed.

### Priority

Critical

---

## EPIC 4 — Timeout Protection

### Tasks

Set:

```
strict timeout
```

Fail safely if provider slow.

### Done Criteria

AI cannot block webhook.

### Priority

Critical

---

## EPIC 5 — AI Audit Logging

### Tasks

Log:

- raw prompt hash
- latency
- confidence
- fallback reason

### Done Criteria

AI behavior inspectable.

### Priority

High

---

## Sprint Exit Criteria

External AI introduced without destabilizing core.

---

## Sprint Risk

Without strict isolation:

AI spreads chaos across architecture.

---

# EPIC 6 — Input / Output Enrichment Layer

## Sprint Goal

Upgrade raw user interaction into clean structured conversational quality before external AI grows.

---

## Tasks

Create:

```
internal/formatter/
```

Files:

```
input_normalizer.go
output_formatter.go
entity_cleaner.go
```

---

## Input Enrichment Responsibilities

Normalize incoming text before parser.

Support:

### punctuation cleanup

```
ingatkan saya!!! besok???
```

---

### shorthand normalization

```
tmrw
pls
jam 8 mlm
```

---

### mixed bilingual input

```
remind aku tomorrow jam 8
```

---

### typo tolerance

```
ingatkn sy besok
```

---

## Output Enrichment Responsibilities

Format canonical result before reply builder.

Support:

### readable time rendering

internal:

```
2026-03-08T20:00:00+07:00
```

user sees:

```
Besok jam 8 malam
```

or:

```
Tomorrow at 8 PM
```

---

### reminder list formatting

```
1. Pay electricity
2. Call mom
3. Renew passport
```

---

### memory retrieval readability

Avoid raw DB-like reply text.

---

### voice transcript cleanup preparation

Voice transcription output must pass same formatter.

---

## Done Criteria

All user-facing text passes formatter layer before parser and before outbound reply.

---

## Priority

Critical

---

## Architecture Rule

Flow becomes:

```
handler
→ formatter input
→ ai parser
→ canonical intent
→ service
→ formatter output
→ whatsapp send
```

---

## Sprint Value

Without formatter isolation:

- parser becomes dirty
- service starts formatting replies
- multilingual quality degrades

---

## Future Unlock

This EPIC enables:

- premium richer replies
- adaptive tone
- summary output
- cleaner voice handling

# Phase 3 Exit Criteria

Phase 3 completes only if:

- recurring reminders stable
- queue active safely
- memory retrievable
- AI fallback controlled

---

# After Phase 3 → Phase 4 Ready

Then NEMORIS can safely evolve into:

- voice notes
- semantic memory ranking
- personal assistant workflows
- calendar sync
- multi-device reliability

---

# Strong Principle for Phase 3

> Add intelligence only where reliability already exists.

Because:

> smart failures are worse than simple reliable systems ⚙️🚀
