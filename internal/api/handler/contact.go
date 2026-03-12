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

// List godoc
// @Summary      List contacts
// @Description  Returns a paginated list of contacts. Supports keyword search on name, phone, and email.
// @Tags         Contacts
// @Produce      json
// @Param        search  query     string  false  "Keyword search"
// @Param        page    query     int     false  "Page number (default 1)"
// @Param        limit   query     int     false  "Page size 1-100 (default 20)"
// @Success      200     {object}  domain.PaginatedResult[domain.Contact]
// @Security     BearerAuth
// @Router       /contacts [get]
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

// Get godoc
// @Summary      Get contact
// @Description  Returns a single contact by UUID.
// @Tags         Contacts
// @Produce      json
// @Param        id  path      string  true  "Contact UUID"
// @Success      200  {object}  domain.Contact
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Security     BearerAuth
// @Router       /contacts/{id} [get]
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

// Create godoc
// @Summary      Create contact
// @Description  Creates a new contact. phone_wa and full_name are required.
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        body  body      domain.Contact  true  "Contact payload"
// @Success      201   {object}  domain.Contact
// @Failure      400   {object}  object{error=string}
// @Security     BearerAuth
// @Router       /contacts [post]
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

// Update godoc
// @Summary      Update contact
// @Description  Partial update — only provided fields are changed.
// @Tags         Contacts
// @Accept       json
// @Produce      json
// @Param        id    path      string          true  "Contact UUID"
// @Param        body  body      domain.Contact  true  "Fields to update"
// @Success      200   {object}  domain.Contact
// @Failure      400   {object}  object{error=string}
// @Failure      404   {object}  object{error=string}
// @Security     BearerAuth
// @Router       /contacts/{id} [patch]
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

// Delete godoc
// @Summary      Delete contact
// @Description  Permanently deletes a contact. Requires admin role.
// @Tags         Contacts
// @Param        id  path  string  true  "Contact UUID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Security     BearerAuth
// @Router       /contacts/{id} [delete]
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
