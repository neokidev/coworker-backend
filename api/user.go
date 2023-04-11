package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	db "github.com/ot07/coworker-backend/db/sqlc"
	"github.com/ot07/coworker-backend/util"
)

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required,email" swaggertype:"string"`
	Password  string `json:"password" validate:"required,min=14"`
}

type userResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required,email" swaggertype:"string"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func (server *Server) createUser(c *fiber.Ctx) error {
	req := new(createUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	arg := db.CreateUserParams{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(c.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.Status(fiber.StatusForbidden).JSON(newErrorResponse(err))
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusOK).JSON(newUserResponse(user))
}
