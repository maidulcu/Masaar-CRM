package api

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/maidulcu/masaar-crm/internal/api/handler"
	"github.com/maidulcu/masaar-crm/internal/api/middleware"
	"github.com/maidulcu/masaar-crm/internal/config"
	"github.com/maidulcu/masaar-crm/internal/ws"
)

type Handlers struct {
	Auth      *handler.AuthHandler
	Contact   *handler.ContactHandler
	Lead      *handler.LeadHandler
	WhatsApp  *handler.WhatsAppHandler
	AI        *handler.AIHandler
}

func RegisterRoutes(app *fiber.App, h *Handlers, hub *ws.Hub, cfg *config.Config) {
	// ── Public routes ────────────────────────────────────────────────────────
	app.Post("/api/v1/auth/login", h.Auth.Login)
	app.Post("/api/v1/auth/refresh", h.Auth.Refresh)

	// WhatsApp webhook — Meta calls this publicly
	app.Get("/webhooks/whatsapp", h.WhatsApp.Verify)
	app.Post("/webhooks/whatsapp", h.WhatsApp.Receive)

	// ── WebSocket — authenticated upgrade ────────────────────────────────────
	app.Use("/ws", func(c *fiber.Ctx) error {
		if fiberws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/warroom", middleware.JWT(cfg.JWTSecret), fiberws.New(hub.Handler()))

	// ── Authenticated API ────────────────────────────────────────────────────
	v1 := app.Group("/api/v1", middleware.JWT(cfg.JWTSecret))

	v1.Delete("/auth/logout", h.Auth.Logout)

	// Contacts
	v1.Get("/contacts", h.Contact.List)
	v1.Post("/contacts", h.Contact.Create)
	v1.Get("/contacts/:id", h.Contact.Get)
	v1.Patch("/contacts/:id", h.Contact.Update)
	v1.Delete("/contacts/:id",
		middleware.RequireRole("admin"),
		h.Contact.Delete,
	)

	// Leads / Pipeline
	v1.Get("/leads", h.Lead.KanbanBoard)
	v1.Post("/leads", h.Lead.Create)
	v1.Get("/leads/:id", h.Lead.Get)
	v1.Patch("/leads/:id/stage", h.Lead.UpdateStage)

	// WhatsApp inbox
	v1.Get("/threads", h.WhatsApp.ListThreads)
	v1.Get("/threads/:id/messages", h.WhatsApp.GetMessages)
	v1.Post("/threads/:id/close", h.WhatsApp.CloseThread)

	// AI
	v1.Post("/ai/score-lead/:id", h.AI.ScoreLead)
	v1.Post("/ai/draft-reply/:thread_id", h.AI.DraftReply)
	v1.Post("/ai/summarize/:thread_id", h.AI.SummarizeThread)

	// Health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
