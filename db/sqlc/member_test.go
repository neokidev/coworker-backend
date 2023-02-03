package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ot07/management-app-demo-backend/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomMember(t *testing.T) Member {
	arg := CreateMemberParams{
		ID:        uuid.New(),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Email:     sql.NullString{String: util.RandomEmail(), Valid: true},
	}

	member, err := testQueries.CreateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, member)

	require.Equal(t, arg.ID, member.ID)
	require.Equal(t, arg.FirstName, member.FirstName)
	require.Equal(t, arg.LastName, member.LastName)
	require.Equal(t, arg.Email, member.Email)

	require.NotZero(t, member.CreatedAt)

	return member
}

func TestCreateMember(t *testing.T) {
	createRandomMember(t)
}

func TestGetMember(t *testing.T) {
	member1 := createRandomMember(t)
	member2, err := testQueries.GetMember(context.Background(), member1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.ID, member2.ID)
	require.Equal(t, member1.FirstName, member2.FirstName)
	require.Equal(t, member1.LastName, member2.LastName)
	require.Equal(t, member1.Email, member2.Email)
	require.WithinDuration(t, member1.CreatedAt.Time, member2.CreatedAt.Time, time.Second)
}

func TestListMember(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomMember(t)
	}

	arg := ListMembersParams{Limit: 5, Offset: 5}

	members, err := testQueries.ListMembers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, members, 5)

	for _, member := range members {
		require.NotEmpty(t, member)
	}
}