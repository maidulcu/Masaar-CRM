package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberws "github.com/gofiber/websocket/v2"
	fiberswagger "github.com/swaggo/fiber-swagger"
	"github.com/maidulcu/masaar-crm/internal/api/handler"
	"github.com/maidulcu/masaar-crm/internal/api/middleware"
	"github.com/maidulcu/masaar-crm/internal/config"
	"github.com/maidulcu/masaar-crm/internal/ws"
	_ "github.com/maidulcu/masaar-crm/docs" // swagger generated docs
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Contact      *handler.ContactHandler
	Lead         *handler.LeadHandler
	WhatsApp     *handler.WhatsAppHandler
	AI           *handler.AIHandler
	Notification *handler.NotificationHandler
	Deal         *handler.DealHandler
	Invoice      *handler.InvoiceHandler
}

// webhookLimiter allows Meta's burst delivery (300 req/min per IP) while
// blocking abuse. Meta retries on 429 so legitimate messages are never lost.
var webhookLimiter = limiter.New(limiter.Config{
	Max:        300,
	Expiration: 1 * time.Minute,
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limit exceeded"})
	},
})

// loginLimiter prevents brute-force on the auth endpoint.
var loginLimiter = limiter.New(limiter.Config{
	Max:        10,
	Expiration: 1 * time.Minute,
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many login attempts"})
	},
})

func RegisterRoutes(app *fiber.App, h *Handlers, hub *ws.Hub, cfg *config.Config) {
	// ── Public routes ────────────────────────────────────────────────────────
	app.Post("/api/v1/auth/login", loginLimiter, h.Auth.Login)
	app.Post("/api/v1/auth/refresh", h.Auth.Refresh)

	// WhatsApp webhook — Meta calls this publicly
	app.Get("/webhooks/whatsapp", h.WhatsApp.Verify)
	app.Post("/webhooks/whatsapp", webhookLimiter, h.WhatsApp.Receive)

	// ── WebSocket — authenticated upgrade ────────────────────────────────────
	app.Use("/ws", func(c *fiber.Ctx) error {
		if fiberws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Personal notifications
	app.Get("/ws/notifications", middleware.JWT(cfg.JWTSecret), fiberws.New(hub.Handler()))

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

	// AI (manual)
	v1.Post("/ai/summarize/:thread_id", h.AI.SummarizeThread)

	// Notifications
	v1.Get("/notifications", h.Notification.List)
	v1.Patch("/notifications/:id/read", h.Notification.MarkRead)

	// Deals
	v1.Get("/deals", h.Deal.List)
	v1.Post("/deals", h.Deal.Create)
	v1.Patch("/deals/:id/stage", h.Deal.UpdateStage)
	v1.Get("/deals/:id/invoices", h.Deal.ListInvoices)

	// Invoices
	v1.Post("/invoices", h.Invoice.Create)
	v1.Get("/invoices/:id", h.Invoice.Get)
	v1.Post("/invoices/:id/send", h.Invoice.Send)
	v1.Patch("/invoices/:id/status", h.Invoice.UpdateStatus)

	// Health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger UI — available in all envs; gate with BasicAuth in production if needed
	app.Get("/docs/*", fiberswagger.WrapHandler)
}
