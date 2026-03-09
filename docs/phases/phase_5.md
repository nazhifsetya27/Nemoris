# **Phase 5 — Commercial Readiness + Market Survival**

> Goal:
>
> Turn NEMORIS into a monetizable, defensible product with sustainable cost structure and operational trust.

This phase should begin only after product durability exists across previous phases.

---

# Sprint 17 — Premium Product Layer

## Sprint Goal

Create monetizable capability without polluting core architecture.

---

## EPIC 1 — Premium Capability Flags

### Tasks

Expand account capability model:

```
free
premium
admin
internal
```

### Done Criteria

Features controlled by capability, not hardcoded checks.

### Priority

Critical

---

## EPIC 2 — Premium Feature Isolation

### Tasks

Gate premium-only features:

- advanced memory retrieval
- recurring complexity
- voice quota increase
- calendar sync depth

### Done Criteria

Premium logic separated cleanly.

### Priority

Critical

---

## EPIC 3 — Soft Limit System

### Tasks

Support free-tier limits:

- reminders/day
- memories/day
- AI fallback/day

### Done Criteria

Limits enforce softly with reply messaging.

### Priority

Critical

---

## EPIC 4 — Upgrade Reply Flow

### Tasks

Localized upgrade prompts.

Example:

```
You reached today's free limit.
```

```
Batas gratis hari ini sudah tercapai.
```

### Done Criteria

Upsell integrated naturally.

### Priority

High

---

## EPIC 5 — Internal Premium Simulation

### Tasks

Run premium internally before public exposure.

### Done Criteria

Premium tested safely.

### Priority

High

---

## Sprint Exit Criteria

Premium exists without architecture contamination.

---

## Sprint Risk

Without isolation:

billing logic spreads everywhere.

---

# Sprint 18 — Cost Control Engine

## Sprint Goal

Prevent growth from silently becoming expensive.

---

## EPIC 1 — Cost Attribution Layer

### Tasks

Track per feature cost:

- WAHA sends
- AI fallback
- transcription
- storage

### Done Criteria

Feature cost visible per user.

### Priority

Critical

---

## EPIC 2 — Heavy Feature Thresholds

### Tasks

Protect expensive features with thresholds.

### Done Criteria

Abnormal usage limited automatically.

### Priority

Critical

---

## EPIC 3 — Cost Dashboard API

### Tasks

Create internal endpoint:

```
/internal/costs
```

### Done Criteria

Daily operational cost visible.

### Priority

High

---

## EPIC 4 — AI Budget Guard

### Tasks

Fallback blocked after quota threshold.

### Done Criteria

AI cannot silently explode cost.

### Priority

Critical

---

## EPIC 5 — Storage Aging Policy

### Tasks

Archive old media/transcription safely.

### Done Criteria

Storage growth controlled.

### Priority

High

---

## Sprint Exit Criteria

Costs measurable before scaling users.

---

## Sprint Risk

Without cost visibility:

growth becomes dangerous.

---

# Sprint 19 — Trust + Retention Layer

## Sprint Goal

Increase user trust so NEMORIS becomes habit, not novelty.

---

## EPIC 1 — Reliability Feedback Replies

### Tasks

After critical reminder delivery:

optional confirmation logic.

### Done Criteria

Trust reinforced subtly.

### Priority

High

---

## EPIC 2 — Reminder History UX

### Tasks

Improve:

```
/reminders
```

presentation quality.

### Done Criteria

Users understand stored state easily.

### Priority

High

---

## EPIC 3 — Memory Recall Quality

### Tasks

Improve recall reply readability.

### Done Criteria

Memory answers feel natural.

### Priority

High

---

## EPIC 4 — User Confidence Signals

### Tasks

Reply examples:

```
Saved successfully.
Reminder scheduled.
Memory updated.
```

### Done Criteria

State becomes explicit.

### Priority

Medium

---

## EPIC 5 — Retention Metrics

### Tasks

Track:

- weekly active users
- reminder repeat usage
- memory reuse

### Done Criteria

Retention measurable.

### Priority

Critical

---

## Sprint Exit Criteria

NEMORIS starts forming habit loops.

---

## Sprint Risk

Without trust:

users churn quietly.

---

# Sprint 20 — Launch Discipline + WhatsApp Survival

## Sprint Goal

Prepare public launch without triggering operational collapse.

---

## EPIC 1 — Controlled User Cohort Launch

### Tasks

Launch by cohort:

```
10 users
50 users
100 users
```

### Done Criteria

Growth staged intentionally.

### Priority

Critical

---

## EPIC 2 — WhatsApp Send Pattern Safety

### Tasks

Strengthen pacing by cohort size.

### Done Criteria

Outbound volume remains human-safe.

### Priority

Critical

---

## EPIC 3 — Ban Risk Monitoring

### Tasks

Track:

- send burst
- repeated failure
- abnormal outbound ratio

### Done Criteria

Ban indicators visible early.

### Priority

Critical

---

## EPIC 4 — Incident Playbook

### Tasks

Prepare:

- WAHA down
- DB unavailable
- queue blocked
- outbound blocked

### Done Criteria

Operational response documented.

### Priority

Critical

---

## EPIC 5 — Founder Launch Dashboard

### Tasks

Single internal dashboard:

- active users
- failures
- due reminders
- queue health
- cost

### Done Criteria

Founder sees business pulse daily.

### Priority

High

---

## Sprint Exit Criteria

NEMORIS can face real users safely.

---

## Sprint Risk

Without launch discipline:

technical success still fails publicly.

---

# Phase 5 Exit Criteria

Phase 5 completes only if:

- premium works
- costs visible
- retention measurable
- launch controlled

---

# After Phase 5 → NEMORIS enters true company phase

Then Phase 6 becomes:

- multi-user shared memory
- team assistant mode
- API platform
- enterprise-grade reliability
- regional expansion

---

# Strong Principle for Phase 5

> Every new paying feature must strengthen trust more than cost.

Because:

> users pay for reliability, not complexity 💰🚀
