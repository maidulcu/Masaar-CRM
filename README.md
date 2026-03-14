<div align="center">

# Masaar CRM — Open-Source WhatsApp CRM for the UAE

**The self-hosted, Arabic-native CRM built for UAE businesses.**

Close deals over WhatsApp · AI thread summaries · Full RTL Arabic UI · PDPL-compliant

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/github/v/release/dynamicweblab/masaar-crm)](https://github.com/dynamicweblab/masaar-crm/releases)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?logo=postgresql)](https://postgresql.org)

**Live Demo:** [masaar.dynamicweblab.com](https://masaar.dynamicweblab.com) &nbsp;·&nbsp; Built by [Dynamic Web Lab FZE LLC](https://dynamicweblab.com)

[Features](#features) · [Quick Start](#quick-start) · [Architecture](#architecture) · [API Reference](#api-reference) · [Contributing](#contributing)

</div>

---

## What is Masaar CRM?

Masaar is a **free, open-source CRM designed specifically for UAE businesses**. Most CRM software was built for Western markets with Arabic added as an afterthought. Masaar is different — WhatsApp-first, Arabic-native, and fully self-hosted so your customer data never leaves the UAE.

**Who is it for?**
- UAE SMEs that close deals over WhatsApp
- Sales teams that need a bilingual (Arabic/English) CRM
- Businesses subject to UAE PDPL data residency requirements
- Developers who want a self-hosted alternative to Salesforce, HubSpot, or Zoho

---

## Features

| Feature | Description |
|---------|-------------|
| **Dashboard Overview** | Real-time stats: total contacts, active leads, open WhatsApp threads, pipeline value, won revenue |
| **WhatsApp Inbox** | Receive and manage inbound WhatsApp messages from Meta Cloud API |
| **Sales Pipeline** | Kanban board with drag-and-drop stage management (New → Won/Lost) |
| **Lead Notes** | Click any pipeline card to view and edit per-lead notes inline |
| **AI Thread Summaries** | One-click AI summarization via local Ollama LLM — data never leaves your server |
| **Deals & VAT Invoicing** | Deal tracking with PDF invoice generation including UAE 5% VAT |
| **Contact Management** | Unified contact profiles linked to WhatsApp numbers and leads |
| **User Settings** | Per-user password change and language preference (Arabic/English) |
| **Role-Based Access** | Three roles: Admin, Agent, Viewer — enforced on every API route |
| **Real-time Notifications** | Personal WebSocket-powered notifications |
| **Full RTL Support** | Complete Arabic interface — layout, typography, and date direction |
| **Audit Log** | Immutable activity log for UAE PDPL compliance |
| **Self-Hosted** | Single Docker binary — deploy to Moro Hub, G42, or any UAE cloud |

---

## Live Demo

Try it at **[masaar.dynamicweblab.com](https://masaar.dynamicweblab.com)**

> Resets every 24 hours. Do not enter real customer data.

---

## Quick Start

**Requirements:** Docker and Docker Compose.

```bash
git clone https://github.com/dynamicweblab/masaar-crm.git
cd masaar-crm
cp .env.example .env
docker compose up
```

| Service | URL |
|---------|-----|
| Dashboard | http://localhost:3000 |
| API | http://localhost:8080/api/v1 |
| Swagger UI | http://localhost:8080/docs |

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
│                 │  WS Hub  │◀── Real-time Notifications  │
│                 └──────────┘                            │
└─────────────────────────────────────────────────────────┘
         ↑ Single Docker binary — no external dependencies
```

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22 + [Fiber](https://gofiber.io) |
| Frontend | Next.js 14 + Tailwind CSS |
| Database | PostgreSQL 16 + pgvector |
| AI | [Ollama](https://ollama.ai) (llama3 / mistral) |
| Real-time | WebSockets (native Fiber) |
| Cache / Sessions | Redis |
| Migrations | [goose](https://github.com/pressly/goose) |
| Auth | JWT (HS256) + bcrypt |

---

## Project Structure

```
masaar-crm/
├── cmd/server/          # Server entry point
├── internal/
│   ├── api/             # HTTP handlers and middleware
│   ├── domain/          # Models and business types
│   ├── repo/            # PostgreSQL repositories
│   ├── ws/              # WebSocket hub
│   └── ai/              # Ollama integration
├── migrations/          # SQL schema (goose)
├── web/
│   ├── app/             # Next.js App Router pages
│   ├── components/      # React UI components
│   ├── store/           # Zustand state management
│   ├── context/         # Language context (AR/EN)
│   └── lib/             # API client, auth helpers
└── docs/                # OpenAPI / Swagger spec
```

---

## API Reference

All protected routes require `Authorization: Bearer <token>`. Full interactive docs at `/docs`.

| Method | Path | Description | Roles |
|--------|------|-------------|-------|
| `POST` | `/api/v1/auth/login` | Authenticate, get token pair | Public |
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
| `POST` | `/api/v1/invoices` | Create VAT invoice | Agent, Admin |
| `POST` | `/api/v1/invoices/:id/send` | Send invoice via WhatsApp | Admin |
| `GET` | `/api/v1/notifications` | List notifications | Auth |
| `PATCH` | `/api/v1/notifications/:id/read` | Mark notification read | Auth |
| `POST` | `/api/v1/ai/summarize/:thread_id` | AI thread summary | Agent, Admin |
| `GET` | `/ws/notifications` | Real-time notification stream | Auth (WS) |

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

Masaar CRM is licensed under the [MIT License](LICENSE). Free to use, modify, and self-host.

---

## Enterprise Edition

Need more for your UAE enterprise? The Enterprise edition adds:

| Feature | Description |
|---------|-------------|
| **Multiple Pipelines** | Sales, HR, Ops — separate pipelines per team |
| **ZATCA E-Invoicing** | UAE FTA compliant with QR code & tax signing |
| **AI Automation** | Auto lead scoring and auto-replies |
| **WhatsApp Sender** | Send outbound messages and media |
| **WhatsApp Bots** | AI-powered chatbots and template messaging |
| **War Room** | Team live leaderboard dashboard |
| **Semantic Search** | AI-powered similarity search across conversations |
| **SSO / SAML** | Enterprise authentication |
| **Emirates ID** | UAE government ID verification |

**Contact:** [dynamicweblab.com](https://dynamicweblab.com) · info@dynamicweblab.com

---

<div align="center">

Built by [Dynamic Web Lab FZE LLC](https://dynamicweblab.com) &nbsp;·&nbsp; UAE 🇦🇪

**Keywords:** open source CRM UAE · WhatsApp CRM · Arabic CRM software · self-hosted CRM · PDPL compliant CRM · CRM for UAE businesses · Go CRM · bilingual CRM Arabic English

مصنوع بعناية للسوق الإماراتي

</div>
