# Changelog

All notable changes to Masaar CRM will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v0.1.0] - 2026-03-13

### Added
- **WhatsApp Receiver** — Receive and read incoming WhatsApp messages via webhook
- **Sales Pipeline** — Single Kanban board with drag-drop stage management
- **AI Summarize** — Manual "Summarize" button powered by Ollama (local LLM)
- **VAT Invoicing** — Simple PDF invoice generator with 5% UAE VAT
- **Notifications** — Personal real-time notifications via WebSocket
- **Keyword Search** — Search contacts by name, phone, email
- **Contact Management** — Unified contact profiles linked to WhatsApp
- **Audit Log** — Immutable activity log for PDPL compliance
- **RTL / LTR** — Full Arabic interface with RTL support
- **Self-hosted** — Single Docker binary deployment

### Tech Stack
- Go 1.22 + Fiber
- Next.js 14 + Tailwind CSS
- PostgreSQL 16
- Ollama (llama3 / mistral)
- WebSockets
- Redis

---

## [Unreleased]

### Features in Progress
- Deals management
- Invoice UI
- Contact/Lead creation forms

### Enterprise (Coming Soon)
- Multiple Pipelines
- ZATCA E-Invoicing
- WhatsApp Sender & Bots
- AI Automation
- Semantic Search
- SSO/SAML
- Emirates ID

---

*For self-hosters: Always check the changelog before upgrading. Breaking changes will be marked with ⚠️.*
