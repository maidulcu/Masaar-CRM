package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/domain"
	"github.com/maidulcu/masaar-crm/internal/repo"
	"github.com/maidulcu/masaar-crm/internal/ws"
)

type LeadHandler struct {
	leads    *repo.LeadRepo
	contacts *repo.ContactRepo
	hub      *ws.Hub
}

func NewLeadHandler(leads *repo.LeadRepo, contacts *repo.ContactRepo, hub *ws.Hub) *LeadHandler {
	return &LeadHandler{leads: leads, contacts: contacts, hub: hub}
}

// GET /api/v1/leads — returns kanban board grouped by stage
func (h *LeadHandler) KanbanBoard(c *fiber.Ctx) error {
	board, err := h.leads.KanbanBoard(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(board)
}

// POST /api/v1/leads
func (h *LeadHandler) Create(c *fiber.Ctx) error {
	var lead domain.Lead
	if err := c.BodyParser(&lead); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if lead.ContactID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "contact_id is required"})
	}
	if lead.Stage == "" {
		lead.Stage = domain.StageNew
	}
	if lead.Currency == "" {
		lead.Currency = "AED"
	}

	if err := h.leads.Create(c.Context(), &lead); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.hub.Broadcast(ws.Event{
		Type:    "lead.created",
		Payload: lead,
	})

	return c.Status(fiber.StatusCreated).JSON(lead)
}

// PATCH /api/v1/leads/:id/stage
func (h *LeadHandler) UpdateStage(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var body struct {
		Stage domain.LeadStage `json:"stage"`
	}
	if err := c.BodyParser(&body); err != nil || body.Stage == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "stage is required"})
	}

	if err := h.leads.UpdateStage(c.Context(), id, body.Stage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.hub.Broadcast(ws.Event{
		Type: "lead.stage_changed",
		Payload: fiber.Map{
			"lead_id": id,
			"stage":   body.Stage,
		},
	})

	return c.JSON(fiber.Map{"lead_id": id, "stage": body.Stage})
}

// GET /api/v1/leads/:id
func (h *LeadHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	lead, err := h.leads.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "lead not found"})
	}

	// Join contact
	contact, _ := h.contacts.GetByID(c.Context(), lead.ContactID)
	lead.Contact = contact

	return c.JSON(lead)
}
