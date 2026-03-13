# Changelog

All notable changes to Masaar CRM will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.1.0] - 2026-03-13

### Added
- **WhatsApp Receiver** — Receive and read incoming WhatsApp messages via Meta webhook
- **Sales Pipeline** — Single Kanban board with drag-drop stage management (dnd-kit)
- **AI Summarize** — Manual "Summarize" button powered by Ollama (local LLM, data never leaves server)
- **Deals Management** — Deal list, stage tracking, linked to leads
- **VAT Invoicing** — Invoice generator with 5% UAE VAT; sequential `INV-YYYY-NNNN` numbering; draft → sent → paid workflow
- **Notifications** — Personal real-time notifications via WebSocket (`/ws/notifications`)
- **Keyword Search** — Search contacts by name, phone, email
- **Contact Management** — Unified contact profiles linked to WhatsApp threads
- **Audit Log** — Immutable activity log for PDPL compliance
- **RTL / LTR** — Full Arabic interface; Cairo font for AR, Inter for EN; logical CSS properties
- **RBAC** — Role-based access control (admin / agent / viewer) enforced on all write routes
- **Rate Limiting** — 300 req/min on WhatsApp webhook; 10 req/min on login
- **OpenAPI / Swagger UI** — Auto-generated spec at `/docs`, `docs/swagger.yaml` committed to repo
- **Self-hosted** — Full stack via `docker compose up` (Next.js + Go + PostgreSQL + Redis + Ollama)

### Tech Stack
- Go 1.22 + Fiber v2
- Next.js 14 + Tailwind CSS (App Router, TypeScript)
- PostgreSQL 16 + pgvector
- Ollama (llama3 / mistral)
- WebSockets (native Fiber)
- Redis (refresh tokens + blacklist)
- goose v3 (SQL migrations)

---

## [Unreleased]

### Enterprise (Coming Soon)
- Multiple Pipelines
- ZATCA E-Invoicing
- WhatsApp Sender & Bots
- AI Automation
- Semantic Search (pgvector)
- SSO/SAML
- Emirates ID

---

*For self-hosters: Always check the changelog before upgrading. Breaking changes will be marked with ⚠️.*
