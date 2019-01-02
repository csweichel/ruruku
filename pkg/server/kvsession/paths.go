package kvsession

import (
	"strings"
)

const pathSeparator = "/"

func pathSession(name string) []byte {
	return []byte(strings.Join([]string{"s", name}, pathSeparator))
}

func pathSessions() []byte {
	return []byte("s")
}

func pathSessionParticipant(session, username string) []byte {
	return []byte(strings.Join([]string{"p", session, username}, pathSeparator))
}

func pathSessionPlan(session string) []byte {
	return []byte(strings.Join([]string{session, "plan"}, pathSeparator))
}
func pathSessionTestcases(session string) []byte {
	return []byte(strings.Join([]string{session, "case"}, pathSeparator))
}

func pathSessionTestcase(session, tc string) []byte {
	return []byte(strings.Join([]string{session, "case", tc}, pathSeparator))
}

func pathSessionClaims(session, tc string) []byte {
	return []byte(strings.Join([]string{session, "claim", tc, ""}, pathSeparator))
}

func pathSessionClaim(session, tc, user string) []byte {
	return []byte(strings.Join([]string{session, "claim", tc, user}, pathSeparator))
}

func pathSessionContributions(session, tc string) []byte {
	return []byte(strings.Join([]string{session, "contrib", tc, ""}, pathSeparator))
}

func pathSessionContribution(session, tc, user string) []byte {
	return []byte(strings.Join([]string{session, "contrib", tc, user}, pathSeparator))
}

func getLastSegment(p []byte) string {
	if len(p) == 0 {
		return ""
	}

	segments := strings.Split(string(p), pathSeparator)
	return segments[len(segments)-1]
}
