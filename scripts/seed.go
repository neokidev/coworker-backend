package main

import (
	"database/sql"
	"fmt"
	"log"

	db "github.com/ot07/coworker-backend/db/sqlc"
	"github.com/ot07/coworker-backend/util"
	"golang.org/x/net/context"

	_ "github.com/lib/pq"
	_ "github.com/ot07/coworker-backend/docs"
)

func setup() (context.Context, *db.SQLStore, error) {
	ctx := context.Background()

	config, err := util.LoadConfig(".")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot load config: %w", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	store := db.NewStore(conn)
	return ctx, store, nil
}

func runSeed(ctx context.Context, store *db.SQLStore) error {
	fmt.Println("Creating user test data...")

	err := db.CreateUserTestData(ctx, store)
	if err != nil {
		return fmt.Errorf("cannot create user test data: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("Starting seed...")

	ctx, store, err := setup()
	if err != nil {
		log.Fatalf("failed to set up: %v", err)
	}

	err = runSeed(ctx, store)
	if err != nil {
		log.Fatalf("failed to run seed: %v", err)
	}

	fmt.Println("Seed completed successfully")
}
