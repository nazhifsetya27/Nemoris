# Nemoris — Architecture & Technical Documentation

## Overview

**Nemoris** is a WhatsApp-native AI memory assistant built to help users store reminders, tasks, notes, and schedules through natural conversation.

Core product goal:

> Users send text, voice, or media via WhatsApp, and Nemoris transforms it into structured memory, reminders, and intelligent follow-up actions.

## Multilingual Support (MVP)

NEMORIS supports **2 languages from MVP stage**: Indonesian (`id`) and English (`en`).

- **Bilingual user input** — Users may send messages in Indonesian or English natural language.
- **Unified internal intent processing** — All messages are converted into one language-neutral internal command model.
- **Localized replies** — Responses are generated in the user's detected language.

> NEMORIS accepts both Indonesian and English user messages while converting them into one language-neutral internal command model.

---

# 1. Product Scope

## Core MVP Features

- WhatsApp text message reminder creation
- Natural language reminder parsing
- Scheduled reminder delivery
- Persistent memory storage
- Task listing
- Voice note transcription
- AI summarization

## Example User Interactions

**English:**

```text
remind me to pay electricity tomorrow 8pm
```

```text
buy milk every monday
```

```text
remember my passport expires in december
```

```text
what are my pending tasks?
```

**Indonesian:**

```text
ingatkan saya untuk bayar listrik besok jam 8 malam
```

Both languages map internally to the same canonical intent.

---

# 2. High-Level Architecture

```mermaid
flowchart TD

A[WhatsApp User] --> B[WAHA WhatsApp Gateway]

B --> C[Webhook API - Go Backend]

C --> D[Message Processor]

D --> E[Language Detection]

E --> F[Intent Parsing]

F --> G[Canonical Internal Intent]

G --> H[Reminder Service]

H --> I[(PostgreSQL)]

H --> J[(Redis Queue)]

J --> K[Worker Scheduler]

K --> L[Localized Reply Builder]

L --> M[Reminder Sender]

M --> B
```

**Message flow with multilingual support:**

```mermaid
flowchart TD
A[Incoming WhatsApp Message] --> B[Language Detection]
B --> C[Intent Parsing]
C --> D[Canonical Internal Intent]
D --> E[Reminder Service]
E --> F[Localized Reply Builder]
```

- Language Detection runs before Intent Parsing.
- Canonical Intent remains language-independent.

---

# 3. Core Architecture Layers

## Messaging Layer

Responsible for WhatsApp communication.

### Technology

- WAHA Free Edition
- Dockerized session management
- Webhook-based event delivery

### Responsibilities

- Session QR login
- Incoming message delivery
- Outgoing text/media sending

---

## Backend Layer

Built in Go.

### Technology

- Go 1.22+
- Gin Framework
- REST API
- JSON Webhooks

### Responsibilities

- Receive webhook events
- Validate incoming payloads
- Route messages
- Trigger AI parsing
- Persist tasks

---

## Intelligence Layer

Transforms human language into structured commands.

### Technology Options

- OpenAI API
- Ollama local model
- Claude API

### Responsibilities

- Intent detection
- Time extraction
- Entity extraction
- Task classification

### Multilingual Contract

The AI parser must return:

- `language` — detected source language (`id` or `en`)
- `intent` — canonical intent (language-independent)
- `task` — extracted task text
- `time` — extracted timestamp

Canonical output is independent of source language.

---

## Persistence Layer

### PostgreSQL

Stores:

- reminders
- tasks
- user memory
- message history

### Redis

Stores:

- delayed jobs
- retry queues
- event scheduling

---

## Worker Layer

Executes background jobs.

### Responsibilities

- Trigger reminders
- Retry failed deliveries
- Run recurring schedules

---

# 4. Container Architecture

```mermaid
flowchart LR

A[Docker Host]

A --> B[WAHA Container]
A --> C[Go API Container]
A --> D[Postgres Container]
A --> E[Redis Container]
A --> F[Worker Container]
```

---

# 5. Docker Compose Design

## Services

| Service  | Purpose            |
| -------- | ------------------ |
| waha     | WhatsApp gateway   |
| backend  | Main API           |
| postgres | Persistent storage |
| redis    | Queue              |
| worker   | Reminder scheduler |

---

# 6. Message Flow

## Incoming Text Message Flow

```mermaid
sequenceDiagram

participant User
participant WAHA
participant API
participant LangDetect
participant Parser
participant DB
participant Queue
participant ReplyBuilder

User->>WAHA: Send WhatsApp message
WAHA->>API: Webhook event
API->>LangDetect: Detect language (id/en)
LangDetect->>Parser: Parse natural language
Parser->>API: Canonical Intent (JSON)
API->>DB: Save reminder
API->>Queue: Schedule task
API->>ReplyBuilder: Build localized reply
ReplyBuilder->>WAHA: Confirmation reply
WAHA->>User: Reminder noted / Pengingat disimpan
```

---

# 7. Reminder Execution Flow

```mermaid
sequenceDiagram

participant Worker
participant Redis
participant DB
participant WAHA
participant User

Worker->>Redis: Pull due job
Worker->>DB: Load reminder
Worker->>WAHA: Send reminder
WAHA->>User: Reminder message
```

---

# 8. Voice Note Flow

## Voice Processing Pipeline

```mermaid
flowchart TD

A[Voice Note Received] --> B[WAHA Media Event]

B --> C[Download Audio]

C --> D[Whisper Transcription]

D --> E[Language Detection]

E --> F[LLM Parser]

F --> G[Canonical Internal Intent]

G --> H[Reminder Service]

H --> I[Localized Reply Builder]
```

---

# 9. AI Parsing Contract

## Input

**English:**

```json
{
  "message": "remind me to pay electricity tomorrow 8pm"
}
```

**Indonesian:**

```json
{
  "message": "ingatkan saya untuk bayar listrik besok jam 8 malam"
}
```

## Output (Canonical Internal Model)

Both inputs map to the same canonical output:

```json
{
  "language": "id",
  "intent": "create_reminder",
  "task": "bayar listrik",
  "time": "2026-03-08T20:00:00+07:00"
}
```

**Critical rule:** `intent` must never depend on human language.

- **Allowed:** `create_reminder`
- **Forbidden:** `buat_pengingat`

## Parser Examples

**English:**

```text
remind me to pay electricity tomorrow 8pm
```

**Indonesian:**

```text
ingatkan saya untuk bayar listrik besok jam 8 malam
```

Both map internally to the same Canonical Intent.

---

# 10. Database Schema

## reminders

```sql
CREATE TABLE reminders (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    task TEXT NOT NULL,
    remind_at TIMESTAMP NOT NULL,
    recurrence TEXT,
    status TEXT DEFAULT 'pending',
    language TEXT DEFAULT 'id',
    created_at TIMESTAMP DEFAULT now()
);
```

- `language` stored for analytics and localized delivery.
- Business logic remains language-neutral.

---

## users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    whatsapp_id TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT now()
);
```

---

## memories

```sql
CREATE TABLE memories (
    id UUID PRIMARY KEY,
    user_id UUID,
    content TEXT,
    type TEXT,
    language TEXT DEFAULT 'id',
    created_at TIMESTAMP DEFAULT now()
);
```

- `language` stored for analytics and localized delivery.

---

# 11. Backend Module Structure

```text
backend/

cmd/api/main.go

internal/
  whatsapp/
  ai/
  reminder/
  scheduler/
  database/
  memory/
  worker/
  service/
    language_detector.go
    reply_builder.go
  i18n/
```

## Internal Multilingual Components

### language_detector.go

- Detects `id` or `en` from incoming message.
- Runs before Intent Parsing.

### reply_builder.go

- Generates response text based on detected language.
- Uses localized templates from `internal/i18n/`.

### internal/i18n/

- Localized message templates.
- Supports Indonesian and English.

---

# 12. Recommended Go Packages

| Concern   | Library        |
| --------- | -------------- |
| HTTP API  | gin            |
| DB ORM    | gorm / sqlx    |
| Queue     | asynq          |
| Migration | golang-migrate |
| Logging   | zap            |
| Config    | viper          |

---

# 13. WAHA Integration Contract

## Send Message

Endpoint:

```text
POST /api/sendText
```

Payload (localized response):

**English:**

```json
{
  "chatId": "628123456789@c.us",
  "text": "Reminder noted"
}
```

**Indonesian:**

```json
{
  "chatId": "628123456789@c.us",
  "text": "Pengingat disimpan"
}
```

---

## Webhook Receive

Expected fields:

```json
{
  "from": "628123456789@c.us",
  "body": "remind me tomorrow"
}
```

---

# 14. Suggested API Endpoints

| Endpoint   | Purpose               |
| ---------- | --------------------- |
| /webhook   | incoming WAHA webhook |
| /health    | health check          |
| /reminders | list reminders        |
| /memory    | list stored memory    |

---

# 15. Scheduling Strategy

## Preferred Approach

Use Redis queue:

```text
Asynq
```

Reason:

- delayed jobs
- retries
- production safe

---

# 16. Production Deployment Architecture

```mermaid
flowchart TD

A[Cloud VPS] --> B[Docker Compose]

B --> C[Nginx]

C --> D[Go API]

D --> E[WAHA]

D --> F[Postgres]

D --> G[Redis]

D --> H[Worker]
```

---

# 17. Security Considerations

## Required

- webhook authentication
- WAHA API token
- encrypted secrets
- database backup

## Recommended

- HTTPS
- reverse proxy
- rate limiting

---

# 18. Future Expansion

## Multilingual Foundation (MVP Implemented)

- **Multilingual foundation implemented in MVP (Indonesian + English)**
- Language Detection before Intent Parsing
- Canonical Intent (language-independent)
- Localized Reply Builder

## Phase 2

- Google Calendar sync
- AI memory retrieval

## Phase 3

- Personal agent mode
- Team memory
- Shared reminders

## Future Multilingual Expansion

- Regional language support
- Voice multilingual parsing

---

# 19. Suggested Repository Naming

```text
nemoris-core
nemoris-api
nemoris-worker
nemoris-ai
nemoris-infra
```

---

# 20. Recommended First Build Order

## Phase 1

1. WAHA connect
2. webhook receive
3. send reply
4. parse reminder
5. save postgres
6. worker reminder

## Phase 2

7. voice transcription
8. recurring tasks
9. memory retrieval

---

# Final Principle

Nemoris should remain:

> chat-first, minimal, invisible infrastructure, strong reliability

Because users only care that:

> they send one message and it always remembers.

---
