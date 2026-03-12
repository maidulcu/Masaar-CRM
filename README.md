<div align="center">

# Masaar CRM

**The open-source CRM built for the UAE market.**

WhatsApp-first pipelines · Agentic AI · Sovereign hosting · Arabic-native

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![Next.js](https://img.shields.io/badge/Next.js-14+-black?logo=next.js)](https://nextjs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?logo=postgresql)](https://postgresql.org)

**Live Demo:** [masaar.dynamicweblab.com](https://masaar.dynamicweblab.com) &nbsp;·&nbsp; Built by [Dynamic Web Lab FZE LLC](https://dynamicweblab.com)

[Features](#features) · [Quick Start](#quick-start) · [Architecture](#architecture) · [Roadmap](#roadmap) · [Contributing](#contributing)

</div>

---

## Why Masaar?

Most CRMs were designed for Western markets and bolted on Arabic support as an afterthought. Masaar is different — built ground-up for the UAE in 2026.

- **WhatsApp is your pipeline.** In the UAE, deals close over WhatsApp. Masaar treats WhatsApp conversations as first-class data, not just an activity log.
- **Your data stays in the UAE.** Ship a single compiled binary to your own server — Moro Hub, G42, or bare metal. No forced cloud, full PDPL compliance.
- **AI that acts, not just reads.** Local LLMs (via Ollama) score your leads and draft context-aware Arabic/English follow-ups — without sending your data to external APIs.
- **Arabic is not a toggle.** The entire dashboard — layout, logic, date formats — flips for Arabic users. RTL-native from day one.

---

## Features

### Core (Open Source — MIT)

| Feature | Description |
|---------|-------------|
| **Sales Pipeline** | Kanban board with drag-drop stage management |
| **WhatsApp Inbox** | Receive, read, and reply to WhatsApp threads in-app |
| **Contact Management** | Unified contact profiles linked to WhatsApp identity |
| **Lead Scoring** | AI-powered scoring via local Ollama LLMs |
| **AI Draft Replies** | Context-aware message drafts in Arabic or English |
| **Live War Room** | Real-time pipeline updates via WebSockets |
| **VAT Invoicing** | Generate and send basic invoices (AED, 5% VAT) |
| **Audit Log** | Immutable activity log for PDPL compliance |
| **RTL / LTR** | Full Arabic interface, no layout hacks |
| **Self-hosted** | Single Docker binary, runs anywhere |

### Enterprise Add-ons

| Feature | Description |
|---------|-------------|
| ZATCA Invoice Signing | UAE FTA e-invoicing with TLV QR codes |
| Advanced AI Agents | Agentic follow-up workflows ("reply to all Burj Khalifa leads") |
| Semantic Search | pgvector-powered contact and thread similarity search |
| Multi-workspace | Manage multiple brands or teams from one instance |
| SSO / SAML | Enterprise authentication for government and finance |
| UAE Cloud Templates | One-click deploy configs for G42 and Moro Hub |

> **Enterprise enquiries:** [dynamicweblab.com](https://dynamicweblab.com) · hello@dynamicweblab.com

---

## Live Demo

A hosted demo is available at **[masaar.dynamicweblab.com](https://masaar.dynamicweblab.com)**.

> The demo resets every 24 hours. Do not enter real customer data.

---

## Quick Start

**Prerequisites:** Docker and Docker Compose installed.

```bash
git clone https://github.com/dynamicweblab/masaar-crm.git
cd masaar-crm
cp .env.example .env
docker compose up
```

The app will be running at:

| Service | URL |
|---------|-----|
| Dashboard | http://localhost:3000 |
| API | http://localhost:8080/api/v1 |
| API Docs | http://localhost:8080/swagger |

Default login: `admin@masaar.local` / `changeme`

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Your Server (UAE)                     │
│                 Moro Hub / G42 / On-prem                │
│                                                          │
│  ┌──────────┐   ┌──────────┐   ┌──────────────────────┐│
│  │ Next.js  │──▶│  Fiber   │──▶│    PostgreSQL         ││
│  │ (RTL/LTR)│   │  (Go)    │   │  + pgvector          ││
│  └──────────┘   │          │──▶│    Redis             ││
│                 │          │   └──────────────────────┘│
│  ┌──────────┐   │          │   ┌──────────────────────┐│
│  │ WhatsApp │──▶│ /webhooks│──▶│  Ollama (local LLM)  ││
│  │ Meta API │   │          │   │  llama3 / mistral    ││
│  └──────────┘   └──────────┘   └──────────────────────┘│
│                      │                                   │
│                 ┌────▼─────┐                            │
│                 │  WS Hub  │◀── Live Sales War Room     │
│                 └──────────┘                            │
└─────────────────────────────────────────────────────────┘
         ↑ Single Docker binary — no external dependencies
```

**Stack:**

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22 + [Fiber](https://gofiber.io) |
| Frontend | Next.js 14 + Tailwind CSS |
| Database | PostgreSQL 16 + pgvector |
| AI | [Ollama](https://ollama.ai) (llama3 / mistral) |
| Real-time | WebSockets (native Fiber) |
| Cache | Redis |
| Migrations | [goose](https://github.com/pressly/goose) |
| Auth | JWT (RS256) |

---

## Project Structure

```
masaar-crm/
├── cmd/server/          # Main entry point
├── internal/
│   ├── api/             # Route handlers
│   ├── domain/          # Business logic & entities
│   ├── repo/            # Database repositories
│   ├── ws/              # WebSocket hub
│   └── ai/              # Ollama AI service
├── migrations/          # SQL migrations
├── web/                 # Next.js frontend
├── docker/              # Dockerfile + Compose
└── docs/                # API documentation
```

---

## Roadmap

- [x] Architecture design and project planning
- [ ] **Week 1** — Core data layer, auth, WhatsApp webhook receiver
- [ ] **Week 2** — Sales pipeline API, Kanban UI, WebSocket war room
- [ ] **Week 3** — AI lead scoring, draft replies, VAT invoicing
- [ ] **Week 4** — RBAC, OpenAPI docs, Docker single-binary, demo

---

## Contributing

Contributions are welcome. Please open an issue first to discuss significant changes.

```bash
# Fork the repo, then:
git clone https://github.com/YOUR_USERNAME/masaar-crm.git
cd masaar-crm
cp .env.example .env
docker compose up -d postgres redis ollama
go run ./cmd/server
```

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting a pull request.

---

## License

Masaar Core is licensed under the [MIT License](LICENSE).

Enterprise modules are proprietary — contact [Dynamic Web Lab FZE LLC](https://dynamicweblab.com) for licensing.

---

<div align="center">

Built by [Dynamic Web Lab FZE LLC](https://dynamicweblab.com) &nbsp;·&nbsp; UAE 🇦🇪

مصنوع بعناية للسوق الإماراتي

</div>
