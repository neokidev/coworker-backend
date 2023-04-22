package db

import (
	"context"

	"github.com/ot07/coworker-backend/util"
)

func CreateUserTestData(ctx context.Context, store *SQLStore) error {
	hashedPassword, err := util.HashPassword("password")
	if err != nil {
		return err
	}

	arg := CreateUserParams{
		FirstName:      "ユーザ",
		LastName:       "テスト",
		Email:          "testuser@email.com",
		HashedPassword: hashedPassword,
	}

	_, err = store.CreateUser(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}
