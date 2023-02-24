package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ot07/coworker-backend/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomMember(t *testing.T, testQueries *Queries) Member {
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

func (s *DatabaseTestSuite) TestCreateMember() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	createRandomMember(t, testQueries)
}

func (s *DatabaseTestSuite) TestGetMember() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	member1 := createRandomMember(t, testQueries)
	member2, err := testQueries.GetMember(context.Background(), member1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member1.ID, member2.ID)
	require.Equal(t, member1.FirstName, member2.FirstName)
	require.Equal(t, member1.LastName, member2.LastName)
	require.Equal(t, member1.Email, member2.Email)
	require.WithinDuration(t, member1.CreatedAt, member2.CreatedAt, time.Second)
}

func (s *DatabaseTestSuite) TestListMembers() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	for i := 0; i < 10; i++ {
		createRandomMember(t, testQueries)
	}

	arg := ListMembersParams{Limit: 5, Offset: 5}

	members, err := testQueries.ListMembers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, members, 5)

	for _, member := range members {
		require.NotEmpty(t, member)
	}
}

func (s *DatabaseTestSuite) TestUpdateMemberAllFields() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	oldMember := createRandomMember(t, testQueries)
	newFirstName := util.RandomName()
	newLastName := util.RandomName()
	newEmail := util.RandomEmail()

	arg := UpdateMemberParams{
		ID:        oldMember.ID,
		FirstName: sql.NullString{String: newFirstName, Valid: true},
		LastName:  sql.NullString{String: newLastName, Valid: true},
		Email:     sql.NullString{String: newEmail, Valid: true},
	}

	updatedMember, err := testQueries.UpdateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedMember)

	require.Equal(t, oldMember.ID, updatedMember.ID)
	require.Equal(t, newFirstName, updatedMember.FirstName)
	require.NotEqual(t, oldMember.FirstName, updatedMember.FirstName)
	require.Equal(t, newLastName, updatedMember.LastName)
	require.NotEqual(t, oldMember.LastName, updatedMember.LastName)
	require.Equal(t, newEmail, updatedMember.Email.String)
	require.NotEqual(t, oldMember.Email.String, updatedMember.Email.String)
	require.WithinDuration(t, oldMember.CreatedAt, updatedMember.CreatedAt, time.Second)
}

func (s *DatabaseTestSuite) TestUpdateMemberOnlyFirstName() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	oldMember := createRandomMember(t, testQueries)
	newFirstName := util.RandomName()

	arg := UpdateMemberParams{
		ID:        oldMember.ID,
		FirstName: sql.NullString{String: newFirstName, Valid: true},
		LastName:  sql.NullString{Valid: false},
		Email:     sql.NullString{Valid: false},
	}

	updatedMember, err := testQueries.UpdateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedMember)

	require.Equal(t, oldMember.ID, updatedMember.ID)
	require.Equal(t, newFirstName, updatedMember.FirstName)
	require.NotEqual(t, oldMember.FirstName, updatedMember.FirstName)
	require.Equal(t, oldMember.LastName, updatedMember.LastName)
	require.Equal(t, oldMember.Email.String, updatedMember.Email.String)
	require.WithinDuration(t, oldMember.CreatedAt, updatedMember.CreatedAt, time.Second)
}

func (s *DatabaseTestSuite) TestUpdateMemberOnlyLastName() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	oldMember := createRandomMember(t, testQueries)
	newLastName := util.RandomName()

	arg := UpdateMemberParams{
		ID:        oldMember.ID,
		FirstName: sql.NullString{Valid: false},
		LastName:  sql.NullString{String: newLastName, Valid: true},
		Email:     sql.NullString{Valid: false},
	}

	updatedMember, err := testQueries.UpdateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedMember)

	require.Equal(t, oldMember.ID, updatedMember.ID)
	require.Equal(t, oldMember.FirstName, updatedMember.FirstName)
	require.Equal(t, newLastName, updatedMember.LastName)
	require.NotEqual(t, oldMember.LastName, updatedMember.LastName)
	require.Equal(t, oldMember.Email.String, updatedMember.Email.String)
	require.WithinDuration(t, oldMember.CreatedAt, updatedMember.CreatedAt, time.Second)
}

func (s *DatabaseTestSuite) TestUpdateMemberOnlyEmail() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	oldMember := createRandomMember(t, testQueries)
	newEmail := util.RandomEmail()

	arg := UpdateMemberParams{
		ID:        oldMember.ID,
		FirstName: sql.NullString{Valid: false},
		LastName:  sql.NullString{Valid: false},
		Email:     sql.NullString{String: newEmail, Valid: true},
	}

	updatedMember, err := testQueries.UpdateMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedMember)

	require.Equal(t, oldMember.ID, updatedMember.ID)
	require.Equal(t, oldMember.FirstName, updatedMember.FirstName)
	require.Equal(t, oldMember.LastName, updatedMember.LastName)
	require.Equal(t, newEmail, updatedMember.Email.String)
	require.NotEqual(t, oldMember.Email.String, updatedMember.Email.String)
	require.WithinDuration(t, oldMember.CreatedAt, updatedMember.CreatedAt, time.Second)
}

func (s *DatabaseTestSuite) TestDeleteMember() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	member1 := createRandomMember(t, testQueries)
	err := testQueries.DeleteMember(context.Background(), member1.ID)
	require.NoError(t, err)

	member2, err := testQueries.GetMember(context.Background(), member1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, member2)
}

func (s *DatabaseTestSuite) TestDeleteMembers() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	member1 := createRandomMember(t, testQueries)
	member2 := createRandomMember(t, testQueries)
	err := testQueries.DeleteMembers(context.Background(), []uuid.UUID{member1.ID, member2.ID})
	require.NoError(t, err)

	member3, err := testQueries.GetMember(context.Background(), member1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, member3)

	member4, err := testQueries.GetMember(context.Background(), member2.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, member4)
}

func (s *DatabaseTestSuite) TestCountMembers() {
	t := s.T()
	t.Parallel()

	tx := beginTransaction(t)
	defer rollbackTransaction(t, tx)

	testQueries := New(tx)

	n := 10
	for i := 0; i < n; i++ {
		createRandomMember(t, testQueries)
	}

	count, err := testQueries.CountMembers(context.Background())
	require.NoError(t, err)
	require.Equal(t, count, int64(n))
}
