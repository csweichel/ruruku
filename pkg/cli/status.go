package cli

import (
    "fmt"
    "github.com/32leaves/ruruku/protocol"
    "github.com/32leaves/ruruku/pkg/storage"
)

type TestCaseResult = protocol.TestCaseResult
const (
    NotRun TestCaseResult = "not-run"
    NotClaimed TestCaseResult = "not-claimed"
)

type TestCaseStatus struct {
    Testcase *protocol.TestCase
    Result protocol.TestCaseResult
    Comments []string
    Claims []string
}

type StatusReport struct {
    Session storage.Storage
    Passed []TestCaseStatus
    Undecided []TestCaseStatus
    Failed []TestCaseStatus
    NotRun []TestCaseStatus
    NotClaimed []*protocol.TestCase

    SuitePassed bool
}

func ComputeStatus(session storage.Storage) StatusReport {
    suite := session.GetSuite()
    run := session.GetRun()

    lists := map[TestCaseResult][]TestCaseStatus{
        protocol.Passed: make([]TestCaseStatus, 0),
        protocol.Undecided: make([]TestCaseStatus, 0),
        protocol.Failed: make([]TestCaseStatus, 0),
    }
    notRun := make([]TestCaseStatus, 0)
    notClaimed := make([]*protocol.TestCase, 0)
    suitePassed := true

    runidx := make(map[string][]*protocol.TestCaseRun)
    for _, tcr := range run.Cases {
        var cse []*protocol.TestCaseRun
        if existingCse, ok := runidx[tcr.CaseID]; ok {
            cse = existingCse
        } else {
            cse = make([]*protocol.TestCaseRun, 0)
        }

        cse = append(cse, &tcr)
        runidx[tcr.CaseID] = cse
    }

    claimidx := make(map[string][]string)
    for _, pcp := range run.Participants {
        for claim, _ := range pcp.ClaimedCases {
            var claims []string
            if existingClaims, ok := claimidx[claim]; ok {
                claims = existingClaims
            } else {
                claims = make([]string, 0)
            }

            claims = append(claims, pcp.Name)
            claimidx[claim] = claims
        }
    }

    for _, tc := range suite.Cases {
        cid := fmt.Sprintf("%s/%s", tc.Group, tc.ID)
        runs, hasRuns := runidx[cid]

        if hasRuns {
            var result = protocol.Passed
            for _, rn := range runs {
                if result == protocol.Passed && rn.Result == protocol.Undecided {
                    result = protocol.Undecided
                    suitePassed = suitePassed && !tc.MustPass
                } else if (result == protocol.Passed || result == protocol.Undecided) && rn.Result == protocol.Failed {
                    result = protocol.Failed
                    suitePassed = suitePassed && !tc.MustPass
                }
            }

            ls, _ := lists[result]
            claims, _ := claimidx[cid]
            lists[result] = append(ls, TestCaseStatus{
                Testcase: &tc,
                Result: result,
                Claims: claims,
            })
        } else {
            claims, hasClaims := claimidx[cid]
            if hasClaims {
                notRun = append(notRun, TestCaseStatus{
                    Testcase: &tc,
                    Result: "claimed",
                    Claims: claims,
                })
            } else {
                notClaimed = append(notClaimed, &tc)
            }

            suitePassed = suitePassed && !tc.MustPass
        }
    }

    return StatusReport{
        Session: session,
        Passed: lists[protocol.Passed],
        Undecided: lists[protocol.Undecided],
        Failed: lists[protocol.Failed],
        NotRun: notRun,
        NotClaimed: notClaimed,
        SuitePassed: suitePassed,
    }
}

func (s *StatusReport) Print(verbose bool) {
    testCount := len(s.Session.GetSuite().Cases)
    fmt.Printf("Suite:   %s\n", s.Session.GetSuite().Name)
    fmt.Printf("Session: %s\n", s.Session.GetRun().Name)
    fmt.Printf("\n")
    fmt.Printf("Passed:      %d/%d", len(s.Passed), testCount)
    if verbose {
        printCases(s.Passed)
    }
    fmt.Printf("\nUndecided:   %d/%d", len(s.Undecided), testCount)
    if verbose {
        printCases(s.Undecided)
    }
    fmt.Printf("\nFailed:      %d/%d", len(s.Failed), testCount)
    if verbose {
        printCases(s.Failed)
    }
    fmt.Printf("\nNot Run:     %d/%d", len(s.NotRun), testCount)
    if verbose {
        printCases(s.NotRun)
    }
    fmt.Printf("\nNot Claimed: %d/%d", len(s.NotClaimed), testCount)
    if len(s.NotClaimed) > 0 {
        fmt.Printf(" (")
    }
    for i, fc := range s.NotClaimed {
        comma := ""
        if i > 0 {
            comma = ", "
        }
        fmt.Printf("%s%s", comma, fc.Name)
    }
    if len(s.NotClaimed) > 0 {
        fmt.Printf(")")
    }
    fmt.Printf("\n\n")
    if s.SuitePassed {
        fmt.Printf("Suite passed\n")
    } else {
        fmt.Printf("Suite failed\n")
    }
}

func printCases(cases []TestCaseStatus) {
    if len(cases) > 0 {
        fmt.Printf(" (")
    }
    for i, fc := range cases {
        comma := ""
        if i > 0 {
            comma = ", "
        }
        fmt.Printf("%s%s", comma, fc.Testcase.Name)
    }
    if len(cases) > 0 {
        fmt.Printf(")")
    }
}