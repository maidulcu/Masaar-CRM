package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maidulcu/masaar-crm/internal/api/middleware"
	"github.com/maidulcu/masaar-crm/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	users *repo.UserRepo
}

func NewUserHandler(users *repo.UserRepo) *UserHandler {
	return &UserHandler{users: users}
}

// GetMe godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's profile.
// @Tags         Users
// @Produce      json
// @Success      200  {object}  object{id=string,name=string,email=string,role=string,lang_pref=string}
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	claims := middleware.ClaimsFromCtx(c)
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	user, err := h.users.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(fiber.Map{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"lang_pref": user.LangPref,
	})
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change the authenticated user's password. Requires current password for verification.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body  object{current_password=string,new_password=string}  true  "Passwords"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Security     BearerAuth
// @Router       /users/me/password [patch]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if body.CurrentPassword == "" || body.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "current_password and new_password are required"})
	}
	if len(body.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new_password must be at least 8 characters"})
	}

	claims := middleware.ClaimsFromCtx(c)
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	user, err := h.users.FindByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "current password is incorrect"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	if err := h.users.UpdatePassword(c.Context(), userID, string(hash)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update password"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateLang godoc
// @Summary      Update language preference
// @Description  Update the authenticated user's language preference (ar or en).
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body  object{lang=string}  true  "Language"
// @Success      200   {object}  object{lang_pref=string}
// @Failure      400   {object}  object{error=string}
// @Security     BearerAuth
// @Router       /users/me/lang [patch]
func (h *UserHandler) UpdateLang(c *fiber.Ctx) error {
	var body struct {
		Lang string `json:"lang"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if body.Lang != "ar" && body.Lang != "en" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "lang must be 'ar' or 'en'"})
	}

	claims := middleware.ClaimsFromCtx(c)
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	if err := h.users.UpdateLangPref(c.Context(), userID, body.Lang); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update language"})
	}

	return c.JSON(fiber.Map{"lang_pref": body.Lang})
}
