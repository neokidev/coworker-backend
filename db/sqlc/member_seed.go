package db

import (
	"context"
	"database/sql"

	"github.com/ot07/coworker-backend/util"
)

func CreateMemberTestData(ctx context.Context, store *SQLStore) error {
	for i := 0; i < 10; i++ {
		arg := CreateMemberParams{
			FirstName: util.RandomName(),
			LastName:  util.RandomName(),
			Email:     sql.NullString{String: util.RandomEmail(), Valid: true},
		}

		_, err := store.CreateMember(ctx, arg)
		if err != nil {
			return err
		}
	}

	return nil
}
