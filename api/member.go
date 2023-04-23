package api

import (
	"database/sql"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/ot07/coworker-backend/db/sqlc"
)

type createMemberRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"omitempty,email" swaggertype:"string" format:"email"`
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
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := db.CreateMemberParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     sql.NullString{String: req.Email, Valid: len(req.Email) > 0},
	}

	member, err := server.store.CreateMember(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type getMemberRequest struct {
	ID uuid.UUID `params:"id"`
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
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	member, err := server.store.GetMember(c.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(newErrorResponse(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type listMembersRequest struct {
	PageID   int32 `query:"page_id" json:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" json:"page_size" validate:"required,min=5,max=10"`
}

type membersResponse []memberResponse

func newMembersResponse(members []db.Member) membersResponse {
	rsp := make(membersResponse, 0, len(members))
	for _, member := range members {
		rsp = append(rsp, newMemberResponse(member))
	}
	return rsp
}

type listMembersResponseMeta struct {
	PageID     int32 `json:"page_id"`
	PageSize   int32 `json:"page_size"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

type listMembersResponse struct {
	Meta listMembersResponseMeta `json:"meta"`
	Data membersResponse         `json:"data"`
}

// @Summary      List members
// @Tags         members
// @Param        query query listMembersRequest true "query"
// @Success      200 {object} listMembersResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members [get]
func (server *Server) listMembers(c *fiber.Ctx) error {
	req := new(listMembersRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
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

	totalCount, err := server.store.CountMembers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	pageCount := int64(math.Ceil(float64(totalCount) / float64(req.PageSize)))

	rsp := listMembersResponse{
		Meta: listMembersResponseMeta{
			PageID:     req.PageID,
			PageSize:   req.PageSize,
			PageCount:  pageCount,
			TotalCount: totalCount,
		},
		Data: newMembersResponse(members),
	}
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type updateMemberRequestParams struct {
	ID uuid.UUID `params:"id" validate:"required"`
}

type updateMemberRequestBody struct {
	FirstName string `json:"first_name" validate:"omitempty" swaggertype:"string"`
	LastName  string `json:"last_name" validate:"omitempty" swaggertype:"string"`
	Email     string `json:"email" validate:"omitempty,email" swaggertype:"string" format:"email"`
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

	body := new(updateMemberRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	arg := db.UpdateMemberParams{
		ID:        params.ID,
		FirstName: sql.NullString{String: body.FirstName, Valid: len(body.FirstName) > 0},
		LastName:  sql.NullString{String: body.LastName, Valid: len(body.LastName) > 0},
		Email:     sql.NullString{String: body.Email, Valid: len(body.Email) > 0},
	}

	member, err := server.store.UpdateMember(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	rsp := newMemberResponse(member)
	return c.Status(fiber.StatusOK).JSON(rsp)
}

type deleteMemberRequest struct {
	ID uuid.UUID `params:"id"`
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
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err := server.store.DeleteMember(c.Context(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}

type deleteMembersRequest struct {
	IDs string `query:"ids" json:"ids" validate:"required"`
}

// @Summary      Delete members
// @Tags         members
// @Param        query query deleteMembersRequest true "query"
// @Success      204 {object} nil
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /members [delete]
func (server *Server) deleteMembers(c *fiber.Ctx) error {
	req := new(deleteMembersRequest)
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	validate := newValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	IDs, err := memberIDsFromCommaSeparatedString(req.IDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(newErrorResponse(err))
	}

	err = server.store.DeleteMembers(c.Context(), IDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(newErrorResponse(err))
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}

func memberIDsFromCommaSeparatedString(commaSeparatedString string) ([]uuid.UUID, error) {
	var IDs []uuid.UUID
	strIDs := strings.Split(commaSeparatedString, ",")
	for _, strID := range strIDs {
		ID, err := uuid.Parse(strID)
		if err != nil {
			return nil, err
		}
		IDs = append(IDs, ID)
	}
	return IDs, nil
}

func memberIDsToCommaSeparatedString(IDs []uuid.UUID) string {
	IDStrings := make([]string, 0, len(IDs))
	for _, ID := range IDs {
		IDStrings = append(IDStrings, ID.String())
	}
	return strings.Join(IDStrings, ",")
}
