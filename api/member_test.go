package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/ot07/coworker-backend/db/mock"
	db "github.com/ot07/coworker-backend/db/sqlc"
	"github.com/ot07/coworker-backend/util"
	"github.com/stretchr/testify/require"
)

func TestGetMemberAPI(t *testing.T) {
	t.Parallel()

	member := randomMember()

	testCases := []struct {
		name          string
		memberID      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:     "OK",
			memberID: member.ID.String(),
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
		{
			name:     "NotFound",
			memberID: member.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(db.Member{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name:     "InternalError",
			memberID: member.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(db.Member{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name:     "InvalidID",
			memberID: "InvalidID",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
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

			// start test server and send request
			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/members/%s", tc.memberID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestCreateMemberAPI(t *testing.T) {
	t.Parallel()

	member := randomMember()
	memberOnlyRequiredFields := db.Member{
		ID:        member.ID,
		FirstName: member.FirstName,
		LastName:  member.LastName,
	}

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
		{
			name: "OptionalFieldsNotFound",
			body: fiber.Map{
				"id":         member.ID,
				"first_name": member.FirstName,
				"last_name":  member.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateMemberParams{
					ID:        member.ID,
					FirstName: member.FirstName,
					LastName:  member.LastName,
				}

				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(memberOnlyRequiredFields, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchMember(t, response.Body, memberOnlyRequiredFields)
			},
		},
		{
			name: "IDNotFound",
			body: fiber.Map{
				"first_name": member.FirstName,
				"last_name":  member.LastName,
				"email":      member.Email.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "FirstNameNotFound",
			body: fiber.Map{
				"id":        member.ID,
				"last_name": member.LastName,
				"email":     member.Email.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "LastNameNotFound",
			body: fiber.Map{
				"id":         member.ID,
				"first_name": member.FirstName,
				"email":      member.Email.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			body: fiber.Map{
				"id":         member.ID,
				"first_name": member.FirstName,
				"last_name":  member.LastName,
				"email":      "InvalidEmail",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "CreateMemberError",
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
					Return(db.Member{}, sql.ErrConnDone)
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

			// start test server and send request
			server := newTestServer(t, store)

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
	t.Parallel()

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

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(1).
					Return(int64(len(members)), nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				checkListMembersResponse(t, response.Body, members, 1, int32(n), 1, int64(n))
			},
		},
		{
			name: "PageIDNotFound",
			query: Query{
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMembers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageIDLessThanLowerLimit",
			query: Query{
				pageID:   0,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMembers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeNotFound",
			query: Query{
				pageID: 1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMembers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeLessThanLowerLimit",
			query: Query{
				pageID:   1,
				pageSize: 4,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMembers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeMoreThanUpperLimit",
			query: Query{
				pageID:   1,
				pageSize: 11,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMembers(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "ListMembersError",
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
					Return([]db.Member{}, sql.ErrConnDone)

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
		{
			name: "CountMembersError",
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

				store.EXPECT().
					CountMembers(gomock.Any()).
					Times(1).
					Return(int64(0), sql.ErrConnDone)
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
	t.Parallel()

	member := randomMember()
	memberOnlyRequiredFields := db.Member{
		ID: member.ID,
	}

	testCases := []struct {
		name          string
		memberID      string
		body          fiber.Map
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:     "OK",
			memberID: member.ID.String(),
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
		{
			name:     "OptionalFieldsNotFound",
			memberID: member.ID.String(),
			body:     fiber.Map{},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateMemberParams{
					ID: member.ID,
				}

				store.EXPECT().
					UpdateMember(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(memberOnlyRequiredFields, nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchMember(t, response.Body, memberOnlyRequiredFields)
			},
		},
		{
			name:     "InvalidEmail",
			memberID: member.ID.String(),
			body: fiber.Map{
				"first_name": member.FirstName,
				"last_name":  member.LastName,
				"email":      "InvalidEmail",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateMember(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:     "UpdateMemberError",
			memberID: member.ID.String(),
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
					Return(db.Member{}, sql.ErrConnDone)
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

			// start test server and send request
			server := newTestServer(t, store)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/members/%s", tc.memberID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestDeleteMemberAPI(t *testing.T) {
	t.Parallel()

	member := randomMember()

	testCases := []struct {
		name          string
		memberID      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name:     "OK",
			memberID: member.ID.String(),
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
		{
			name:     "InvalidID",
			memberID: "InvalidID",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:     "DeleteMemberError",
			memberID: member.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMember(gomock.Any(), gomock.Eq(member.ID)).
					Times(1).
					Return(sql.ErrConnDone)
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

			// start test server and send request
			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/members/%s", tc.memberID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			tc.checkResponse(t, response)
		})
	}
}

func TestDeleteMembersAPI(t *testing.T) {
	t.Parallel()

	member1 := randomMember()
	member2 := randomMember()
	memberIDs := []uuid.UUID{member1.ID, member2.ID}

	type Query struct {
		IDs string
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
				IDs: memberIDsToCommaSeparatedString(memberIDs),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMembers(gomock.Any(), gomock.Eq(memberIDs)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNoContent, response.StatusCode)
			},
		},
		{
			name:  "IDsNotFound",
			query: Query{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMembers(gomock.Any(), gomock.Eq(memberIDs)).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InvalidIDs",
			query: Query{
				IDs: memberIDsToCommaSeparatedString(memberIDs) + ",InvalidID",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMembers(gomock.Any(), gomock.Eq(memberIDs)).
					Times(0)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
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

			// start test server and send request
			server := newTestServer(t, store)

			url := "/api/v1/members"
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("ids", fmt.Sprintf("%s", tc.query.IDs))
			request.URL.RawQuery = q.Encode()

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

func checkListMembersResponse(t *testing.T, body io.ReadCloser, members []db.Member, pageID int32, pageSize int32, pageCount int64, totalCount int64) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse listMembersResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)

	require.Equal(t, pageID, gotResponse.Meta.PageID)
	require.Equal(t, pageSize, gotResponse.Meta.PageSize)
	require.Equal(t, pageCount, gotResponse.Meta.PageCount)
	require.Equal(t, totalCount, gotResponse.Meta.TotalCount)

	gotMembers := gotResponse.Data
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
