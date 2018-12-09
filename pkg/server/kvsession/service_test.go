package kvsession

//go:generate go run ../../../build/generate-api-tests.go kvsession $GOFILE

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	bolt "github.com/etcd-io/bbolt"
	"os"
)

func newTestServer() api.SessionServiceServer {
	if _, err := os.Stat("test.db"); err == nil {
		os.Remove("test.db")
	}

	db, err := bolt.Open("test.db", 0666, nil)
	if err != nil {
		panic(err)
	}

	store, err := NewSession(db)
	if err != nil {
		panic(err)
	}

	return store
}
