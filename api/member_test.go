package api

import (
	"encoding/json"
	"fmt"
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

			url := fmt.Sprintf("/members/%s", member.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
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
	}
}

func requireBodyMatchMember(t *testing.T, body io.ReadCloser, member db.Member) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMember db.Member
	err = json.Unmarshal(data, &gotMember)
	require.NoError(t, err)
	require.Equal(t, member, gotMember)

	err = body.Close()
	require.NoError(t, err)
}
