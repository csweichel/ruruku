package types

import (
	"fmt"
	"strings"
)

type ParticipantToken struct {
	SessionID     string
	ParticipantID string
}

func ParseParticipantToken(token string) (*ParticipantToken, error) {
	segments := strings.Split(token, "/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Invalid participant token")
	}

	sid := strings.Trim(segments[0], " ")
	if sid == "" {
		return nil, fmt.Errorf("Invalid participant token")
	}
	pid := segments[1]
	if pid == "" {
		return nil, fmt.Errorf("Invalid participant token")
	}

	return &ParticipantToken{SessionID: sid, ParticipantID: pid}, nil
}

func (t *ParticipantToken) String() string {
	return fmt.Sprintf("%s/%s", t.SessionID, t.ParticipantID)
}
