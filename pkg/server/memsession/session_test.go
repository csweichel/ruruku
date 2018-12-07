package memsession

//go:generate go run ../../../build/generate-api-tests.go memsession $GOFILE

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
)

func newTestServer() api.SessionServiceServer {
	return NewMemoryBackedSessionStore()
}
