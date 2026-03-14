<div align="center">

# Masaar CRM

**The open-source CRM built for the UAE market.**

WhatsApp-first CRM · AI Summarize · Arabic-native · Self-hosted

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/github/v/release/dynamicweblab/masaar-crm)](https://github.com/dynamicweblab/masaar-crm/releases)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?logo=postgresql)](https://postgresql.org)

**Live Demo:** [masaar.dynamicweblab.com](https://masaar.dynamicweblab.com) &nbsp;·&nbsp; Built by [Dynamic Web Lab FZE LLC](https://dynamicweblab.com)

[Features](#features) · [Quick Start](#quick-start) · [Architecture](#architecture) · [Roadmap](#roadmap) · [Contributing](#contributing)

</div>

---

## Why Masaar?

Most CRMs were designed for Western markets and bolted on Arabic support as an afterthought. Masaar is different — built ground-up for the UAE in 2026.

- **WhatsApp is your pipeline.** In the UAE, deals close over WhatsApp. Masaar treats WhatsApp conversations as first-class data, not just an activity log.
- **Your data stays in the UAE.** Ship a single compiled binary to your own server — Moro Hub, G42, or bare metal. No forced cloud, full PDPL compliance.
- **AI that summarizes.** Manual "Summarize" button for threads — local LLMs via Ollama, your data never leaves your server.
- **Arabic is not a toggle.** The entire dashboard — layout, logic, date formats — flips for Arabic users. RTL-native from day one.

---

## Features

| Feature | Description |
|---------|-------------|
| **Dashboard Overview** | Real-time stat cards: contacts, active leads, open threads, open deals, won revenue |
| **WhatsApp Receiver** | Receive and read incoming WhatsApp messages |
| **Sales Pipeline** | Kanban board with drag-drop stage management and inline lead notes |
| **Lead Notes** | Click any pipeline card to view/edit per-lead notes without leaving the board |
| **AI Summarize** | Manual "Summarize" button for threads |
| **VAT Invoicing** | Simple PDF invoice generator with 5% VAT |
| **Notifications** | Personal real-time notifications via WebSocket |
| **Keyword Search** | Standard keyword search for contacts |
| **Contact Management** | Unified contact profiles linked to WhatsApp |
| **User Settings** | Change password (bcrypt-verified) and language preference per user |
| **Audit Log** | Immutable activity log for PDPL compliance |
| **RTL / LTR** | Full Arabic interface, persisted per user |
| **Self-hosted** | Single Docker binary, runs anywhere |

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
| API Docs (Swagger UI) | http://localhost:8080/docs |
| OpenAPI Spec | `docs/swagger.yaml` |

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
│  │ (RTL/LTR)│   │  (Go)    │   │                      ││
│  └──────────┘   │          │──▶│    Redis             ││
│                 │          │   └──────────────────────┘│
│  ┌──────────┐   │          │   ┌──────────────────────┐│
│  │ WhatsApp │──▶│ /webhooks│──▶│  Ollama (local LLM)  ││
│  │ Meta API │   │          │   │  llama3 / mistral    ││
│  └──────────┘   └──────────┘   └──────────────────────┘│
│                      │                                   │
│                 ┌────▼─────┐                            │
│                 │  WS Hub  │◀── Personal Notifications   │
│                 └──────────┘                            │
└─────────────────────────────────────────────────────────┘
         ↑ Single Docker binary — no external dependencies
```

**Stack:**

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22 + [Fiber](https://gofiber.io) |
| Frontend | Next.js 14 + Tailwind CSS |
| Database | PostgreSQL 16 |
| AI | [Ollama](https://ollama.ai) (llama3 / mistral) |
| Real-time | WebSockets (native Fiber) |
| Cache | Redis |
| Migrations | [goose](https://github.com/pressly/goose) |
| Auth | JWT (HS256) + bcrypt |

---

## API Reference

All protected routes require `Authorization: Bearer <token>`.

| Method | Path | Description | Roles |
|--------|------|-------------|-------|
| `POST` | `/api/v1/auth/login` | Login, returns token pair | Public |
| `POST` | `/api/v1/auth/refresh` | Refresh access token | Public |
| `DELETE` | `/api/v1/auth/logout` | Invalidate session | Auth |
| `GET` | `/api/v1/stats` | Dashboard overview metrics | Auth |
| `GET` | `/api/v1/users/me` | Current user profile | Auth |
| `PATCH` | `/api/v1/users/me/password` | Change password | Auth |
| `PATCH` | `/api/v1/users/me/lang` | Update language preference | Auth |
| `GET` | `/api/v1/contacts` | List contacts (search + paginate) | Auth |
| `POST` | `/api/v1/contacts` | Create contact | Agent, Admin |
| `PATCH` | `/api/v1/contacts/:id` | Update contact | Agent, Admin |
| `DELETE` | `/api/v1/contacts/:id` | Delete contact | Admin |
| `GET` | `/api/v1/leads` | Kanban board (all stages) | Auth |
| `POST` | `/api/v1/leads` | Create lead | Agent, Admin |
| `PATCH` | `/api/v1/leads/:id/stage` | Move lead to stage | Agent, Admin |
| `PATCH` | `/api/v1/leads/:id/notes` | Update lead notes | Agent, Admin |
| `GET` | `/api/v1/threads` | List WhatsApp threads | Auth |
| `GET` | `/api/v1/threads/:id/messages` | Thread messages | Auth |
| `POST` | `/api/v1/threads/:id/close` | Close thread | Agent, Admin |
| `GET` | `/api/v1/deals` | List deals | Auth |
| `POST` | `/api/v1/deals` | Create deal | Agent, Admin |
| `PATCH` | `/api/v1/deals/:id/stage` | Update deal stage | Agent, Admin |
| `GET` | `/api/v1/invoices/:id` | Get invoice | Auth |
| `POST` | `/api/v1/invoices` | Create invoice | Agent, Admin |
| `POST` | `/api/v1/invoices/:id/send` | Send invoice via WhatsApp | Admin |
| `GET` | `/api/v1/notifications` | List notifications | Auth |
| `PATCH` | `/api/v1/notifications/:id/read` | Mark notification read | Auth |
| `POST` | `/api/v1/ai/summarize/:thread_id` | AI thread summary | Agent, Admin |
| `GET` | `/ws/notifications` | Personal notification stream | Auth (WS) |

Full interactive docs at `/docs` (Swagger UI).

---

## Project Structure

```
masaar-crm/
├── cmd/server/          # Main entry point
├── internal/
│   ├── api/             # Route handlers and middleware
│   ├── domain/          # Models and business types
│   ├── repo/            # Database repositories
│   ├── ws/              # WebSocket hub
│   └── ai/              # Ollama AI service
├── migrations/          # SQL migrations (goose)
├── web/
│   ├── app/             # Next.js App Router pages
│   ├── components/      # React UI components
│   ├── store/           # Zustand state (auth)
│   ├── context/         # Language context (AR/EN)
│   └── lib/             # API client, auth helpers
├── docker/              # Dockerfile + Compose
└── docs/                # OpenAPI / Swagger spec
```

---

## Roadmap

- [x] Architecture design and project planning
- [x] **Week 1** — Core data layer, auth, WhatsApp webhook receiver
- [x] **Week 2** — Next.js frontend: Kanban pipeline, WhatsApp inbox, contacts, notifications
- [x] **Week 3** — Deals management, VAT invoicing, full-stack Docker Compose
- [x] **Week 4** — RBAC on all routes, OpenAPI/Swagger UI, rate limiting, v0.1.0 release
- [x] **Week 5** — User settings (change password, language), lead notes, dashboard stats overview

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

Masaar is licensed under the [MIT License](LICENSE).

---

## Enterprise Edition

Need more? The Enterprise edition adds:

| Feature | Description |
|---------|-------------|
| **Multiple Pipelines** | Sales, HR, Ops — separate pipelines |
| **ZATCA E-Invoicing** | UAE FTA compliant with QR & tax signing |
| **AI Automation** | Auto lead scoring & auto-replies |
| **WhatsApp Sender** | Send outbound messages & media |
| **WhatsApp Bots** | AI-powered chatbots & template messaging |
| **War Room** | Team live leaderboard dashboard |
| **Semantic Search** | AI-powered similarity search |
| **SSO / SAML** | Enterprise authentication |
| **Emirates ID** | Government ID verification |

**Contact:** [dynamicweblab.com](https://dynamicweblab.com) · info@dynamicweblab.com

---

<div align="center">

Built by [Dynamic Web Lab FZE LLC](https://dynamicweblab.com) &nbsp;·&nbsp; UAE 🇦🇪

مصنوع بعناية للسوق الإماراتي

</div>
