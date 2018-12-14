package kvsession

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	bolt "github.com/etcd-io/bbolt"
    log "github.com/sirupsen/logrus"
	"os"
)

func newTestServer() (api.SessionServiceServer, api.UserServiceServer) {
    log.SetLevel(log.WarnLevel)

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

	return store, nil
}
