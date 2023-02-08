package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	db "github.com/ot07/coworker-backend/db/sqlc"
)

var validate = validator.New()

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

	app.Post("/members", server.createMember)
	app.Get("/members/:id", server.getMember)
	app.Get("/members", server.listMembers)
	app.Put("/members/:id", server.updateMember)
	app.Delete("/members/:id", server.deleteMember)

	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}

func errorResponse(err error) fiber.Map {
	return fiber.Map{"error": err.Error()}
}
