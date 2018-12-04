package types

import "testing"

func TestParseToken(t *testing.T) {
	token, err := ParseParticipantToken("foo/bar")
	if err != nil {
		t.Errorf("ParseParticipantToken returned error despite valid input: %v", err)
	}
	if token.SessionID != "foo" {
		t.Errorf("ParseParticipantToken did not parse session ID correctly: %v != foo", token.SessionID)
	}
	if token.ParticipantID != "bar" {
		t.Errorf("ParseParticipantToken did not parse participant ID correctly: %v != bar", token.ParticipantID)
	}

	token, err = ParseParticipantToken("/bar")
	if err == nil {
		t.Errorf("ParseParticipantToken did not return an error despite empty sessionID")
	}
	token, err = ParseParticipantToken("bar/")
	if err == nil {
		t.Errorf("ParseParticipantToken did not return an error despite empty participantID")
	}
	token, err = ParseParticipantToken("invalid-session-token")
	if err == nil {
		t.Errorf("ParseParticipantToken did not return an error when parsing an invalid token")
	}
	token, err = ParseParticipantToken("")
	if err == nil {
		t.Errorf("ParseParticipantToken did not return an error when parsing an invalid token")
	}
}

func TestTokenString(t *testing.T) {
	original := ParticipantToken{SessionID: "foo", ParticipantID: "bar"}
	str := original.String()
	token, err := ParseParticipantToken(str)
	if err != nil {
		t.Errorf("ParseParticipantToken(token.String()) returned an error")
	}
	if token.SessionID != original.SessionID {
		t.Errorf("ParseParticipantToken(token.String()) does not yield correct session ID: %s != %s", original.SessionID, token.SessionID)
	}
	if token.ParticipantID != original.ParticipantID {
		t.Errorf("ParseParticipantToken(token.String()) does not yield correct participant ID: %s != %s", original.ParticipantID, token.ParticipantID)
	}
}
