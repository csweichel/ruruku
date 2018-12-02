package server

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	apitests "github.com/32leaves/ruruku/pkg/api/v1/test"
	log "github.com/sirupsen/logrus"
	"testing"
)

var implementations = map[string]func() api.SessionServiceServer{
	"session_mem": func() api.SessionServiceServer { return NewMemoryBackedSessionStore() },
}

func init() {
	log.SetLevel(log.WarnLevel)
}

func TestStartValidSession(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestStartValidSession(t, srv())
		})
	}
}

func TestStartInvalidSession(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestStartInvalidSession(t, srv())
		})
	}
}

func TestTestValidRegistration(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestValidRegistration(t, srv())
		})
	}
}

func TestDuplicateRegistration(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestDuplicateRegistration(t, srv())
		})
	}
}

func TestInvalidRegistration(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestInvalidRegistration(t, srv())
		})
	}
}

func TestBasicStatus(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestBasicStatus(t, srv())
		})
	}
}

func TestValidClaim(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestValidClaim(t, srv())
		})
	}
}

func TestInvalidClaim(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestInvalidClaim(t, srv())
		})
	}
}

func TestValidContribution(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestValidContribution(t, srv())
		})
	}
}

func TestInvalidContribution(t *testing.T) {
	for nm, srv := range implementations {
		t.Run(nm, func(t *testing.T) {
			apitests.RunTestInvalidContribution(t, srv())
		})
	}
}
