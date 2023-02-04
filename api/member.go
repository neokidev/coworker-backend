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

type memberResponse struct {
	ID        uuid.UUID         `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     db.NullString     `json:"email"`
	Status    db.MemberStatuses `json:"status"`
	CreatedAt db.NullTime       `json:"created_at"`
}

func newMemberResponse(user db.Member) memberResponse {
	return memberResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     db.NullString{NullString: user.Email},
		Status:    user.Status,
		CreatedAt: db.NullTime{NullTime: user.CreatedAt},
	}
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

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}
