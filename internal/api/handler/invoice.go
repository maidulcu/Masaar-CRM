package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/domain"
	"github.com/maidulcu/masaar-crm/internal/repo"
)

type InvoiceHandler struct {
	invoices *repo.InvoiceRepo
	deals    *repo.DealRepo
}

func NewInvoiceHandler(invoices *repo.InvoiceRepo, deals *repo.DealRepo) *InvoiceHandler {
	return &InvoiceHandler{invoices: invoices, deals: deals}
}

// POST /api/v1/invoices
func (h *InvoiceHandler) Create(c *fiber.Ctx) error {
	var body struct {
		DealID   uuid.UUID `json:"deal_id"`
		Subtotal float64   `json:"subtotal"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if body.DealID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "deal_id is required"})
	}
	if body.Subtotal <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "subtotal must be positive"})
	}

	invoiceNo, err := h.invoices.NextInvoiceNo(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	inv := domain.VATInvoice{
		DealID:    body.DealID,
		InvoiceNo: invoiceNo,
		Subtotal:  body.Subtotal,
		VATRate:   0.05,
		Status:    domain.InvoiceDraft,
	}

	if err := h.invoices.Create(c.Context(), &inv); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(inv)
}

// GET /api/v1/invoices/:id
func (h *InvoiceHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	inv, err := h.invoices.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invoice not found"})
	}
	return c.JSON(inv)
}

// POST /api/v1/invoices/:id/send
func (h *InvoiceHandler) Send(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.invoices.UpdateStatus(c.Context(), id, domain.InvoiceSent); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id, "status": domain.InvoiceSent})
}

// PATCH /api/v1/invoices/:id/status
func (h *InvoiceHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var body struct {
		Status domain.InvoiceStatus `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil || body.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "status required"})
	}
	if err := h.invoices.UpdateStatus(c.Context(), id, body.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id, "status": body.Status})
}
