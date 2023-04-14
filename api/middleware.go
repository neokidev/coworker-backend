package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ot07/coworker-backend/token"
)

const (
	accessTokenCookieKey  = "access_token"
	accessTokenTypeBearer = "bearer"
	accessTokenPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessTokenWithType := c.Cookies(accessTokenCookieKey)
		if len(accessTokenWithType) == 0 {
			err := errors.New("access token not found")
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		fields := strings.Fields(accessTokenWithType)
		if len(fields) < 2 {
			err := errors.New("invalid access token format")
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		accessTokenType := strings.ToLower(fields[0])
		if accessTokenType != accessTokenTypeBearer {
			err := fmt.Errorf("unsupported access token type %s", accessTokenType)
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(newErrorResponse(err))
		}

		c.Locals(accessTokenPayloadKey, payload)
		return c.Next()
	}
}
