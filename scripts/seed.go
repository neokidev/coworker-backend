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

func truncateAllTables(ctx context.Context, store *db.SQLStore) error {
	log.Println("truncating all tables...")

	err := store.TruncateMembersTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate members table:", err)
	}

	err = store.TruncateSessionsTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate sessions table:", err)
	}

	err = store.TruncateUsersTable(ctx)
	if err != nil {
		log.Fatal("cannot truncate users table:", err)
	}

	return nil
}

func runSeed(ctx context.Context, store *db.SQLStore) error {
	log.Println("creating user test data...")
	err := db.CreateUserTestData(ctx, store)
	if err != nil {
		log.Fatal("cannot create user test data:", err)
	}

	log.Println("creating member test data...")
	err = db.CreateMemberTestData(ctx, store)
	if err != nil {
		log.Fatal("cannot create member test data:", err)
	}

	return nil
}

func main() {
	log.Println("starting seed...")

	ctx, store, err := setup()
	if err != nil {
		log.Fatalf("failed to set up: %v", err)
	}

	err = truncateAllTables(ctx, store)
	if err != nil {
		log.Fatalf("failed to truncate all tables: %v", err)
	}

	err = runSeed(ctx, store)
	if err != nil {
		log.Fatalf("failed to run seed: %v", err)
	}

	log.Println("seed completed successfully")
}
