package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/api/middleware"
	"github.com/maidulcu/masaar-crm/internal/domain"
	"github.com/maidulcu/masaar-crm/internal/repo"
)

type DealHandler struct {
	deals    *repo.DealRepo
	invoices *repo.InvoiceRepo
}

func NewDealHandler(deals *repo.DealRepo, invoices *repo.InvoiceRepo) *DealHandler {
	return &DealHandler{deals: deals, invoices: invoices}
}

// GET /api/v1/deals
func (h *DealHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}

	result, err := h.deals.List(c.Context(), nil, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

// GET /api/v1/deals/:id
func (h *DealHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	deal, err := h.deals.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "deal not found"})
	}
	return c.JSON(deal)
}

// POST /api/v1/deals
func (h *DealHandler) Create(c *fiber.Ctx) error {
	var deal domain.Deal
	if err := c.BodyParser(&deal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if deal.LeadID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "lead_id is required"})
	}
	if deal.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	if deal.Currency == "" {
		deal.Currency = "AED"
	}
	if deal.Stage == "" {
		deal.Stage = domain.DealStageOpen
	}
	if deal.Probability == 0 {
		deal.Probability = 50
	}

	// Set owner from JWT
	if sub, ok := middleware.ClaimsFromCtx(c)["sub"].(string); ok {
		deal.OwnerID, _ = uuid.Parse(sub)
	}

	if err := h.deals.Create(c.Context(), &deal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(deal)
}

// PATCH /api/v1/deals/:id/stage
func (h *DealHandler) UpdateStage(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var body struct {
		Stage domain.DealStage `json:"stage"`
	}
	if err := c.BodyParser(&body); err != nil || body.Stage == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "stage required"})
	}
	if err := h.deals.UpdateStage(c.Context(), id, body.Stage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"id": id, "stage": body.Stage})
}

// GET /api/v1/deals/:id/invoices
func (h *DealHandler) ListInvoices(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	invoices, err := h.invoices.ListByDeal(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(invoices)
}

