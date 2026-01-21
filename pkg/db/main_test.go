package db_test

import (
	"codim/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

// assertJSONEqual compares two JSON RawMessage values by unmarshaling and comparing
func assertJSONEqual(t *testing.T, expected *json.RawMessage, actual *json.RawMessage, fieldName string) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		require.Fail(t, fmt.Sprintf("%s: one is nil, other is not", fieldName))
		return
	}

	var expectedVal interface{}
	var actualVal interface{}

	err := json.Unmarshal(*expected, &expectedVal)
	require.NoError(t, err, fmt.Sprintf("%s: failed to unmarshal expected", fieldName))

	err = json.Unmarshal(*actual, &actualVal)
	require.NoError(t, err, fmt.Sprintf("%s: failed to unmarshal actual", fieldName))

	require.Equal(t, expectedVal, actualVal, fmt.Sprintf("%s: JSON values don't match", fieldName))
}
