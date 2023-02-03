package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/ot07/management-app-demo-backend/db/sqlc"
)

type createMemberRequest struct {
	ID        uuid.UUID     `json:"id" validate:"required"`
	FirstName string        `json:"first_name" validate:"required"`
	LastName  string        `json:"last_name" validate:"required"`
	Email     db.NullString `json:"email" validate:"email"`
}

func (server *Server) createMember(c *fiber.Ctx) error {
	req := new(createMemberRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.CreateMemberParams{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email.NullString,
	}

	member, err := server.store.CreateMember(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return c.Status(fiber.StatusOK).JSON(member)
}
