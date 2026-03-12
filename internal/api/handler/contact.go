package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/domain"
	"github.com/maidulcu/masaar-crm/internal/repo"
)

type ContactHandler struct {
	contacts *repo.ContactRepo
}

func NewContactHandler(contacts *repo.ContactRepo) *ContactHandler {
	return &ContactHandler{contacts: contacts}
}

// GET /api/v1/contacts
func (h *ContactHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := h.contacts.List(c.Context(), search, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

// GET /api/v1/contacts/:id
func (h *ContactHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	contact, err := h.contacts.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact not found"})
	}
	return c.JSON(contact)
}

// POST /api/v1/contacts
func (h *ContactHandler) Create(c *fiber.Ctx) error {
	var contact domain.Contact
	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if contact.PhoneWA == "" || contact.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "phone_wa and full_name are required"})
	}
	if contact.Language == "" {
		contact.Language = "en"
	}

	if err := h.contacts.Create(c.Context(), &contact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(contact)
}

// PATCH /api/v1/contacts/:id
func (h *ContactHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	existing, err := h.contacts.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact not found"})
	}

	// Merge only provided fields
	var patch domain.Contact
	if err := c.BodyParser(&patch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if patch.FullName != "" {
		existing.FullName = patch.FullName
	}
	if patch.Email != "" {
		existing.Email = patch.Email
	}
	if patch.Language != "" {
		existing.Language = patch.Language
	}
	if patch.LeadScore != 0 {
		existing.LeadScore = patch.LeadScore
	}
	if patch.AssignedTo != nil {
		existing.AssignedTo = patch.AssignedTo
	}

	if err := h.contacts.Update(c.Context(), existing); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(existing)
}

// DELETE /api/v1/contacts/:id
func (h *ContactHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.contacts.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
