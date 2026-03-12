package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/config"
	"github.com/maidulcu/masaar-crm/internal/domain"
	"github.com/maidulcu/masaar-crm/internal/repo"
	"github.com/maidulcu/masaar-crm/internal/ws"
)

type WhatsAppHandler struct {
	wa       *repo.WhatsAppRepo
	contacts *repo.ContactRepo
	hub      *ws.Hub
	config   *config.Config
}

func NewWhatsAppHandler(wa *repo.WhatsAppRepo, contacts *repo.ContactRepo, hub *ws.Hub, cfg *config.Config) *WhatsAppHandler {
	return &WhatsAppHandler{wa: wa, contacts: contacts, hub: hub, config: cfg}
}

// GET /webhooks/whatsapp — Meta webhook verification
func (h *WhatsAppHandler) Verify(c *fiber.Ctx) error {
	mode      := c.Query("hub.mode")
	token     := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.config.WAVerifyToken {
		return c.SendString(challenge)
	}
	return c.SendStatus(fiber.StatusForbidden)
}

// GET /api/v1/threads
func (h *WhatsAppHandler) ListThreads(c *fiber.Ctx) error {
	status := c.Query("status", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	threads, err := h.wa.ListThreads(c.Context(), status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(threads)
}

// GET /api/v1/threads/:id/messages
func (h *WhatsAppHandler) GetMessages(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	limit, _ := strconv.Atoi(c.Query("limit", "100"))

	msgs, err := h.wa.GetMessages(c.Context(), id, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(msgs)
}

// POST /api/v1/threads/:id/close
func (h *WhatsAppHandler) CloseThread(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.wa.CloseThread(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// POST /webhooks/whatsapp — receive inbound messages
func (h *WhatsAppHandler) Receive(c *fiber.Ctx) error {
	var payload struct {
		Object string `json:"object"`
		Entry  []struct {
			ID      string `json:"id"`
			Changes []struct {
				Value struct {
					MessagingProduct string `json:"messaging_product"`
					Metadata         struct {
						PhoneNumberID string `json:"phone_number_id"`
					} `json:"metadata"`
					Contacts []struct {
						Profile struct {
							Name string `json:"name"`
						} `json:"profile"`
						WAID string `json:"wa_id"`
					} `json:"contacts"`
					Messages []struct {
						From      string `json:"from"`
						ID        string `json:"id"`
						Timestamp string `json:"timestamp"`
						Type      string `json:"type"`
						Text      struct {
							Body string `json:"body"`
						} `json:"text"`
					} `json:"messages"`
				} `json:"value"`
				Field string `json:"field"`
			} `json:"changes"`
		} `json:"entry"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusOK) // always 200 to Meta
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}
			val := change.Value

			for _, msg := range val.Messages {
				if msg.Type != "text" {
					continue
				}

				senderName := msg.From
				for _, wc := range val.Contacts {
					if wc.WAID == msg.From {
						senderName = wc.Profile.Name
						break
					}
				}

				contact, err := h.contacts.Upsert(c.Context(), msg.From, senderName)
				if err != nil {
					log.Printf("whatsapp: upsert contact error: %v", err)
					continue
				}

				thread, err := h.wa.UpsertThread(c.Context(), contact.ID, val.Metadata.PhoneNumberID)
				if err != nil {
					log.Printf("whatsapp: upsert thread error: %v", err)
					continue
				}

				waMsg := &domain.WhatsAppMessage{
					ThreadID:    thread.ID,
					Direction:   domain.DirectionInbound,
					Body:        msg.Text.Body,
					WAMessageID: msg.ID,
				}
				if err := h.wa.SaveMessage(c.Context(), waMsg); err != nil {
					log.Printf("whatsapp: save message error: %v", err)
					continue
				}

				_ = h.wa.UpdateThreadMeta(c.Context(), thread.ID)

				h.hub.Broadcast(ws.Event{
					Type: "whatsapp.message",
					Payload: fiber.Map{
						"thread_id": thread.ID,
						"contact":   contact,
						"message":   waMsg,
					},
				})
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
