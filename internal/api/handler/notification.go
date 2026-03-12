package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/api/middleware"
	"github.com/maidulcu/masaar-crm/internal/repo"
)

type NotificationHandler struct {
	notifications *repo.NotificationRepo
}

func NewNotificationHandler(notifications *repo.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{notifications: notifications}
}

func userIDFromCtx(c *fiber.Ctx) (uuid.UUID, bool) {
	claims := middleware.ClaimsFromCtx(c)
	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(sub)
	return id, err == nil
}

// GET /api/v1/notifications
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userID, ok := userIDFromCtx(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	notifications, err := h.notifications.ListByUser(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(notifications)
}

// PATCH /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	userID, ok := userIDFromCtx(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.notifications.MarkRead(c.Context(), id, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
