package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ot07/coworker-backend/token"
	"github.com/ot07/coworker-backend/util"
	"github.com/stretchr/testify/require"
)

func addAccessTokenInCookie(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	accessTokenType string,
	userID uuid.UUID,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)

	cookie := &http.Cookie{
		Name:     accessTokenCookieKey,
		Value:    fmt.Sprintf("%s %s", accessTokenType, token),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	request.AddCookie(cookie)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAccessTokenInCookie(t, request, tokenMaker, accessTokenTypeBearer, util.RandomUUID(), time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAccessTokenInCookie(t, request, tokenMaker, "unsupported", util.RandomUUID(), time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAccessTokenInCookie(t, request, tokenMaker, "", util.RandomUUID(), time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAccessTokenInCookie(t, request, tokenMaker, accessTokenTypeBearer, util.RandomUUID(), -time.Minute)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.app.Get(
				authPath,
				authMiddleware(server.tokenMaker),
				func(c *fiber.Ctx) error {
					return c.SendStatus(fiber.StatusOK)
				},
			)

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)

		})
	}
}
