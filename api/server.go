package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	db "github.com/ot07/coworker-backend/db/sqlc"
)

// Server serves HTTP requests for this app service.
type Server struct {
	store db.Store
	app   *fiber.App
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	app := fiber.New()
	app.Use(cors.New())

	server := &Server{store: store, app: app}

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/users", server.createUser)

	v1.Post("/members", server.createMember)
	v1.Get("/members/:id", server.getMember)
	v1.Get("/members", server.listMembers)
	v1.Put("/members/:id", server.updateMember)
	v1.Delete("/members/:id", server.deleteMember)
	v1.Delete("/members", server.deleteMembers)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(err error) errorResponse {
	return errorResponse{Error: err.Error()}
}
