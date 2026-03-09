# NEMORIS Docs — Read on Your Phone

*Last updated: 2026-03-07*

## Quick Start

```bash
cd docs
node serve.js
```

Then open the printed URL on your phone (e.g. `http://192.168.1.x:3847`).

---

## Access Options

### 1. Same WiFi (local network)

Your phone and Mac must be on the same WiFi.

- Run `node serve.js` in the `docs` folder
- Note the "On your phone" URL (e.g. `http://192.168.1.5:3847`)
- Open that URL in your phone's browser

### 2. Public access (from anywhere)

Use **ngrok** to get a public URL:

```bash
# Install ngrok: brew install ngrok
# Start the docs server first (node serve.js in another terminal)
ngrok http 3847
```

ngrok will show a public URL like `https://abc123.ngrok.io` — open it on your phone from anywhere.

---

## What You Get

- **Index page** — Choose which doc to read
- **Progress Summary** — Build status & phases (includes multilingual MVP)
- **Master Roadmap** — Product evolution
- **Project Instructions** — Dev setup & principles
- **Architecture** — Technical docs (Language Detection, Canonical Intent, Localized Reply)

All docs render as readable HTML on mobile.
