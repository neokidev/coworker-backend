package api

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	mockdb "github.com/ot07/coworker-backend/db/mock"
	db "github.com/ot07/coworker-backend/db/sqlc"
	"github.com/ot07/coworker-backend/util"
	"github.com/stretchr/testify/require"
)

func addSessionTokenInCookie(
	request *http.Request,
	sessionToken string,
) {
	cookie := &http.Cookie{
		Name:     sessionTokenKey,
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	request.AddCookie(cookie)
}

func buildValidSessionStubs(store *mockdb.MockStore, session db.Session) {
	store.EXPECT().
		GetSession(gomock.Any(), gomock.Eq(session.SessionToken)).
		Times(1).
		Return(session, nil)
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	session := randomSession()
	expiredSession := randomExpiredSession()

	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStubs: func(store *mockdb.MockStore) {
				buildValidSessionStubs(store, session)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(request *http.Request) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InvalidTokenFormat",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, "invalid")
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, expiredSession.SessionToken.String())
			},
			buildStubs: func(store *mockdb.MockStore) {
				buildValidSessionStubs(store, expiredSession)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(request *http.Request) {
				addSessionTokenInCookie(request, session.SessionToken.String())
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.SessionToken)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			authPath := "/auth"
			server.app.Get(
				authPath,
				authMiddleware(server),
				func(c *fiber.Ctx) error {
					return c.SendStatus(fiber.StatusOK)
				},
			)

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(request)
			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func randomSession() db.Session {
	return db.Session{
		ID:           util.RandomUUID(),
		UserID:       util.RandomUUID(),
		SessionToken: util.RandomUUID(),
		ExpiredAt:    time.Now().Add(time.Minute),
	}
}

func randomExpiredSession() db.Session {
	return db.Session{
		ID:           util.RandomUUID(),
		UserID:       util.RandomUUID(),
		SessionToken: util.RandomUUID(),
		ExpiredAt:    time.Now().Add(-time.Minute),
	}
}
