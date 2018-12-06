package kvsession

//go:generate go run ../../../build/generate-api-tests.go kvsession $GOFILE

import (
    "os"
    api "github.com/32leaves/ruruku/pkg/api/v1"
)

func newTestServer() api.SessionServiceServer {
    if _, err := os.Stat("test.db"); err == nil {
        os.Remove("test.db")
    }

    store, err := NewSession("test.db")
    if err != nil {
        panic(err)
    }

    return store
}
