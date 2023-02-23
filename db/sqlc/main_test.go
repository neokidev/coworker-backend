package db

import (
	"database/sql"
	"github.com/ot07/coworker-backend/util"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

type DatabaseTestSuite struct {
	suite.Suite
}

func TestDatabaseTestSuite(t *testing.T) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	suite.Run(t, new(DatabaseTestSuite))
}
