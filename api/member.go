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

func newMemberResponse(member db.Member) memberResponse {
	return memberResponse{
		ID:        member.ID,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     db.NullString{NullString: member.Email},
		Status:    member.Status,
		CreatedAt: db.NullTime{NullTime: member.CreatedAt},
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

type getMemberRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

func (server *Server) getMember(c *fiber.Ctx) error {
	req := new(getMemberRequest)

	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	member, err := server.store.GetMember(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type listMembersRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

type membersResponse []memberResponse

func newMembersResponse(members []db.Member) membersResponse {
	var rsp []memberResponse
	for _, member := range members {
		rsp = append(rsp, newMemberResponse(member))
	}
	return rsp
}

func (server *Server) listMembers(c *fiber.Ctx) error {
	req := new(listMembersRequest)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.ListMembersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	members, err := server.store.ListMembers(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	rsp := newMembersResponse(members)
	return c.Status(fiber.StatusOK).JSON(rsp)
}
