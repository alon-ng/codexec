package db_test

import (
	"codim/pkg/db"
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	config, err := db.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := db.NewPool(context.Background(), config)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(connPool)
	os.Exit(m.Run())
}

func getRandomInt() int {
	return rand.Intn(100000)
}
