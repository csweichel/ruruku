package storage

import "github.com/32leaves/ruruku/protocol"

type Storage interface {
    ClaimTestcase(tester string, caseID string, claim bool) error
    SetTestcaseRun(tester string, caseID string, result protocol.TestCaseResult, comment string) error
    AddParticipant(tester string) error
    Save() error

    GetSuite() protocol.TestSuite
    GetRun() protocol.TestRun
    GetParticipant(tester string) (protocol.TestParticipant, bool)
}
