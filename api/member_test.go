package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/ot07/coworker-backend/db/mock"
	db "github.com/ot07/coworker-backend/db/sqlc"
	"github.com/ot07/coworker-backend/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestGetMemberAPI(t *testing.T) {
	member := randomMember()

	testCases := []struct {
		name          string
		memberID      uuid.UUID
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:     "OK",
			memberID: member.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(member, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchMember(t, response.Body, member)
			},
		},
		// TODO: add more cases
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)

			url := fmt.Sprintf("/api/v1/members/%s", member.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestCreateMemberAPI(t *testing.T) {
	member := randomMember()

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"id":         member.ID,
				"first_name": member.FirstName,
				"last_name":  member.LastName,
				"email":      member.Email.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateMemberParams{
					ID:        member.ID,
					FirstName: member.FirstName,
					LastName:  member.LastName,
					Email:     member.Email,
				}

				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(member, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchMember(t, response.Body, member)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/members"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestListMembersAPI(t *testing.T) {
	n := 5
	members := make([]db.Member, n)
	for i := 0; i < n; i++ {
		members[i] = randomMember()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListMembersParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListMembers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(members, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchMembers(t, response.Body, members)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)

			url := "/api/v1/members"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestUpdateMemberAPI(t *testing.T) {
	member := randomMember()

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"first_name": member.FirstName,
				"last_name":  member.LastName,
				"email":      member.Email.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateMemberParams{
					ID:        member.ID,
					FirstName: sql.NullString{String: member.FirstName, Valid: true},
					LastName:  sql.NullString{String: member.LastName, Valid: true},
					Email:     member.Email,
				}

				store.EXPECT().
					UpdateMember(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(member, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchMember(t, response.Body, member)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/members/%s", member.ID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestDeleteMemberAPI(t *testing.T) {
	member := randomMember()

	testCases := []struct {
		name          string
		memberID      uuid.UUID
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:     "OK",
			memberID: member.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
			},
		},
		// TODO: add more cases
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)

			url := fmt.Sprintf("/api/v1/members/%s", member.ID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func randomMember() db.Member {
	return db.Member{
		ID:        util.RandomUUID(),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     sql.NullString{String: util.RandomEmail(), Valid: true},
	}
}

func requireBodyMatchMember(t *testing.T, body io.ReadCloser, member db.Member) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMember memberResponse
	err = json.Unmarshal(data, &gotMember)
	require.NoError(t, err)

	requireMemberResponseMatchMember(t, gotMember, member)

	err = body.Close()
	require.NoError(t, err)
}

func requireBodyMatchMembers(t *testing.T, body io.ReadCloser, members []db.Member) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMembers []memberResponse
	err = json.Unmarshal(data, &gotMembers)
	require.NoError(t, err)

	require.Equal(t, len(members), len(gotMembers))
	for i := 0; i < len(members); i++ {
		requireMemberResponseMatchMember(t, gotMembers[i], members[i])
	}

	err = body.Close()
	require.NoError(t, err)
}

func requireMemberResponseMatchMember(t *testing.T, gotMember memberResponse, member db.Member) {
	require.Equal(t, member.ID, gotMember.ID)
	require.Equal(t, member.FirstName, gotMember.FirstName)
	require.Equal(t, member.LastName, gotMember.LastName)
	require.Equal(t, member.Email.String, gotMember.Email.String)
	require.Equal(t, member.CreatedAt, gotMember.CreatedAt)
}
