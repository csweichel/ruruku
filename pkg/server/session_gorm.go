package server

import (
	"context"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/jinzhu/gorm"
)

type gormBackedSessionStore struct {
	DB       *gorm.DB
	Notifier *Notifier
}

type session struct {
	ID     string
	Name   string
	PlanID string
	Open   bool

    Participants []participant
    Cases        []testcaseStatus
}

type participant struct {
    gorm.Model

    session *session `gorm:"foreignkey:sessionID"`
    sessionID string
	types.Participant
}

type testcaseStatus struct {
	sessionID string
	types.Testcase

    Claims  []participant `gorm:"many2many:testcase_participants;"`
    Results []testcaseResult
}

type testcaseResult struct {
    testcaseStatusID string

	Participant   participant
	types.TestcaseRunResult
}

func NewGormBackedSessionStore(db *gorm.DB) *gormBackedSessionStore {
	db.AutoMigrate(
		&session{},
		&participant{},
		&testcaseStatus{},
		&testcaseResult{},
	)

	return &gormBackedSessionStore{
		DB:       db,
		Notifier: NewNotifier(),
	}
}

// Version returns the service release name, service version.
func (s *gormBackedSessionStore) Version(context.Context, *api.VersionRequest) (*api.VersionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Start begins a new test run session
func (s *gormBackedSessionStore) Start(context.Context, *api.StartSessionRequest) (*api.StartSessionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Close closes a session so that no further participation is possible
func (s *gormBackedSessionStore) Close(context.Context, *api.CloseSessionRequest) (*api.CloseSessionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// List returns all sessions
func (s *gormBackedSessionStore) List(*api.ListSessionsRequest, api.SessionService_ListServer) error {
	return fmt.Errorf("Not implemented")
}

// Register adds a new participant to the session. Returns an error
// if a participant with that name already exists.
func (s *gormBackedSessionStore) Register(context.Context, *api.RegistrationRequest) (*api.RegistrationResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Claim expresses a participants intent to execute a testcase
func (s *gormBackedSessionStore) Claim(context.Context, *api.ClaimRequest) (*api.ClaimResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Contribute submitts the result of a testcas execution/run.
// The same participant may contribute to the same testcase multiple times,
// which updates previous contributions.
func (s *gormBackedSessionStore) Contribute(context.Context, *api.ContributionRequest) (*api.ContributionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Status returns the status of a test run
func (s *gormBackedSessionStore) Status(context.Context, *api.SessionStatusRequest) (*api.SessionStatusResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Subscribes to updates on a test run
func (s *gormBackedSessionStore) Updates(*api.SessionUpdatesRequest, api.SessionService_UpdatesServer) error {
	return fmt.Errorf("Not implemented")
}
