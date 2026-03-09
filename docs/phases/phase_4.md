# **Phase 4 — Product Scale + Market Readiness**

> Goal:
>
> Transform NEMORIS into a market-capable WhatsApp-native personal memory system:
>
> voice, semantic memory, external integrations, and production-grade operational durability.

This phase directly extends the future roadmap already implied in your architecture documents: voice transcription, semantic memory, Google Calendar sync, and stronger deployment maturity
architecture

---

# Sprint 13 — Voice Intelligence Layer

## Sprint Goal

Allow users to interact naturally through voice notes.

---

## EPIC 1 — Voice Intake Pipeline

### Tasks

Detect WAHA voice/media payload.

Create:

```
internal/voice/
```

Files:

```
extract.go
normalize.go
```

### Done Criteria

Voice message safely extracted.

### Priority

Critical

---

## EPIC 2 — Audio Storage Safety

### Tasks

Temporary storage policy:

- save file short-term
- delete after processing

### Done Criteria

Audio does not accumulate.

### Priority

Critical

---

## EPIC 3 — Transcription Adapter

### Tasks

Create provider adapter:

```
internal/voice/provider/
```

### Done Criteria

Transcription isolated from business logic.

### Priority

Critical

---

## EPIC 4 — Voice → Existing Intent Pipeline

### Tasks

Flow:

```
voice
→ transcription
→ ai parser
→ canonical intent
```

### Done Criteria

Voice uses same parser as text.

### Priority

Critical

---

## EPIC 5 — Voice Failure Handling

### Tasks

If transcription fails:

localized fallback reply.

### Done Criteria

Failure visible to user safely.

### Priority

High

---

## Sprint Exit Criteria

Voice note creates reminders/memory safely.

---

## Sprint Risk

Without strict pipeline:

voice introduces instability fast.

---

# Sprint 14 — Semantic Memory Layer

## Sprint Goal

Upgrade memory from keyword lookup into relevance-based retrieval.

---

## EPIC 1 — Memory Search Scoring

### Tasks

Rank memory by:

- keyword proximity
- recency
- sender scope

### Done Criteria

Most relevant memory returned first.

### Priority

Critical

---

## EPIC 2 — Semantic Preparation Layer

### Tasks

Create:

```
internal/memory/indexer.go
```

### Done Criteria

Memory prepared for later embeddings.

### Priority

Critical

---

## EPIC 3 — Multi-Match Clarification Engine

### Tasks

If similar results:

ask clarification.

### Done Criteria

Wrong memory avoided.

### Priority

High

---

## EPIC 4 — Retrieval Confidence Layer

### Tasks

Return confidence score internally.

### Done Criteria

Low-confidence retrieval visible.

### Priority

High

---

## EPIC 5 — Memory Aging Strategy

### Tasks

Support:

```
recent
important
archived
```

### Done Criteria

Memory lifecycle begins.

### Priority

Medium

---

## Sprint Exit Criteria

Memory feels intelligent, not raw database search.

---

## Sprint Risk

Without ranking:

memory becomes frustrating quickly.

---

# Sprint 15 — External Calendar + Time System

## Sprint Goal

Allow reminders to connect with external schedule systems.

---

## EPIC 1 — Calendar Integration Adapter

### Tasks

Create:

```
internal/integration/calendar/
```

### Done Criteria

Calendar isolated from reminder core.

### Priority

Critical

---

## EPIC 2 — Google Calendar Sync MVP

### Tasks

Support:

- create event
- update event

### Done Criteria

Reminder optionally mirrored externally.

### Priority

Critical

---

## EPIC 3 — Timezone Safety Layer

### Tasks

Explicit timezone normalization.

### Done Criteria

External sync keeps correct local time.

### Priority

Critical

---

## EPIC 4 — Sync Failure Isolation

### Tasks

Calendar failure must not break reminder creation.

### Done Criteria

Reminder always survives sync issue.

### Priority

High

---

## EPIC 5 — Future Integration Contract

### Tasks

Prepare adapter for:

- Apple Calendar
- Outlook

### Done Criteria

No rewrite needed later.

### Priority

Medium

---

## Sprint Exit Criteria

NEMORIS can extend outside WhatsApp safely.

---

## Sprint Risk

Without isolation:

integration pollutes core architecture.

---

# Sprint 16 — Production Durability + Monetization Readiness

## Sprint Goal

Prepare NEMORIS for real public users and business model decisions.

---

## EPIC 1 — User Tier Foundation

### Tasks

Add account capability flags:

```
free
premium
internal
```

### Done Criteria

Future monetization possible.

### Priority

Critical

---

## EPIC 2 — Usage Metering

### Tasks

Track:

- reminders created
- memory saves
- voice usage
- AI fallback calls

### Done Criteria

Operational cost visible per user.

### Priority

Critical

---

## EPIC 3 — Anti-Abuse Protection Upgrade

### Tasks

Strengthen:

- burst detection
- suspicious repetition
- abnormal command loops

### Done Criteria

Public abuse harder.

### Priority

Critical

---

## EPIC 4 — Backup + Restore Verification

### Tasks

Test:

```
full restore
```

from production dump.

### Done Criteria

Recovery proven.

### Priority

Critical

---

## EPIC 5 — Deployment Hardening

### Tasks

Split:

- staging
- production

Add:

- secrets discipline
- rollback safety

### Done Criteria

Deployment becomes safe.

### Priority

Critical

---

## Sprint Exit Criteria

System survives public usage pressure.

---

## Sprint Risk

Without this sprint:

growth breaks trust fast.

---

# Phase 4 Exit Criteria

Phase 4 completes only if:

- voice stable
- semantic memory useful
- external sync safe
- production durability proven

---

# After Phase 4 → NEMORIS becomes real product territory

Then Phase 5 can begin:

- premium assistant mode
- proactive memory
- shared memory spaces
- business API
- enterprise-grade reliability

---

# Strong Principle for Phase 4

> Features must now justify operational cost.

Because:

> at scale, every feature becomes infrastructure 💡⚙️
