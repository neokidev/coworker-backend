package api

import (
	_ "database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/ot07/coworker-backend/db/sqlc"
	"time"
)

type createMemberRequest struct {
	ID        uuid.UUID     `json:"id" validate:"required" format:"uuid"`
	FirstName string        `json:"first_name" validate:"required"`
	LastName  string        `json:"last_name" validate:"required"`
	Email     db.NullString `json:"email" validate:"email" swaggertype:"string" format:"email"`
}

type memberResponse struct {
	ID        uuid.UUID     `json:"id"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Email     db.NullString `json:"email" swaggertype:"string"`
	CreatedAt time.Time     `json:"created_at"`
}

func newMemberResponse(member db.Member) memberResponse {
	return memberResponse{
		ID:        member.ID,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     db.NullString{NullString: member.Email},
		CreatedAt: member.CreatedAt,
	}
}

// @Summary      Create member
// @Tags         members
// @Param        body body createMemberRequest true "Member object"
// @Success      200 {object} memberResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members [post]
func (server *Server) createMember(c *fiber.Ctx) error {
	req := new(createMemberRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := db.CreateMemberParams{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email.NullString,
	}

	member, err := server.store.CreateMember(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type getMemberRequest struct {
	ID uuid.UUID `params:"id" validate:"required"`
}

// @Summary      Get member
// @Tags         members
// @Param        id path string true "Member ID"
// @Success      200 {object} memberResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members/{id} [get]
func (server *Server) getMember(c *fiber.Ctx) error {
	req := new(getMemberRequest)

	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	member, err := server.store.GetMember(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
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

// @Summary      List members
// @Tags         members
// @Param        query query listMembersRequest true "query"
// @Success      200 {object} membersResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members [get]
func (server *Server) listMembers(c *fiber.Ctx) error {
	req := new(listMembersRequest)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := db.ListMembersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	members, err := server.store.ListMembers(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMembersResponse(members)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type updateMemberRequestParams struct {
	ID uuid.UUID `params:"id" validate:"required"`
}

type updateMemberRequestBody struct {
	FirstName db.NullString `json:"first_name" swaggertype:"string"`
	LastName  db.NullString `json:"last_name" swaggertype:"string"`
	Email     db.NullString `json:"email" validate:"email" swaggertype:"string" format:"email"`
}

// @Summary      Update member
// @Tags         members
// @Param        id   path string                  true "Member ID"
// @Param        body body updateMemberRequestBody true "Member object"
// @Success      200 {object} memberResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members/{id} [put]
func (server *Server) updateMember(c *fiber.Ctx) error {
	params := new(updateMemberRequestParams)

	if err := c.ParamsParser(params); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	if err := validate.Struct(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	body := new(updateMemberRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := db.UpdateMemberParams{
		ID:        params.ID,
		FirstName: body.FirstName.NullString,
		LastName:  body.LastName.NullString,
		Email:     body.Email.NullString,
	}

	member, err := server.store.UpdateMember(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type deleteMemberRequest struct {
	ID uuid.UUID `params:"id" validate:"required"`
}

// @Summary      Delete member
// @Tags         members
// @Param        id path string true "Member ID"
// @Success      204 {object} nil
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members/{id} [delete]
func (server *Server) deleteMember(c *fiber.Ctx) error {
	req := new(deleteMemberRequest)

	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err := server.store.DeleteMember(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
