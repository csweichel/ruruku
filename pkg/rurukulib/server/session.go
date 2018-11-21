package server

import (
	"encoding/json"
	"fmt"
	"github.com/32leaves/ruruku/protocol"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type sessionStore struct {
	suite        *protocol.TestSuite
	run          *protocol.TestRun
	Name         string
	conn         map[*websocket.Conn]string
	participants map[string]*protocol.TestParticipant
}

func LoadSessionOrNew(name string, suite *protocol.TestSuite) (*sessionStore, error) {
    r := &sessionStore{
		suite: suite,
		run: &protocol.TestRun{
			SuiteName: suite.Name,
			Start:     time.Now().Format(time.RFC3339),
		},
		Name:         name,
		conn:         make(map[*websocket.Conn]string),
		participants: make(map[string]*protocol.TestParticipant),
	}

	if _, err := os.Stat(name); err == nil {
		fc, err := ioutil.ReadFile(name)
		if err != nil {
			log.WithField("session", name).WithError(err).Error("Cannot read session file")
			return nil, err
		}

        if err := yaml.Unmarshal(fc, r.run); err != nil {
			log.WithField("session", name).WithError(err).Error("Error while restoring session")
			return nil, err
		}

        for _, pcp := range r.run.Participants {
            r.participants[pcp.Name] = &pcp
        }

		log.WithField("session", name).Info("Restored session")
	}

	return r, nil
}

func (session *sessionStore) Join(conn *websocket.Conn) error {
	return nil
}

func (session *sessionStore) HandleMessage(conn *websocket.Conn, msg []byte) error {
	var rawMsg map[string]interface{}
	if err := json.Unmarshal(msg, &rawMsg); err != nil {
		return err
	}

	if tpe, ok := rawMsg["type"]; !ok {
		return fmt.Errorf("Invalid message does not have a type: %s", msg)
	} else if tpe == "welcome" {
		welcome, err := protocol.UnmarshalWelcomeRequest(msg)
		if err != nil {
			return err
		}
		return session.handleWelcome(conn, welcome)
	} else if tpe == "claim" {
		claim, err := protocol.UnmarshalClaimRequest(msg)
		if err != nil {
			return err
		}
		return session.handleClaim(conn, claim)
    } else if tpe == "newTestCaseRun" {
        newRun, err := protocol.UnmarshalNewTestCaseRunRequest(msg)
        if err != nil {
            return err
        }
        return session.handleNewTestcaseRun(conn, newRun)
	} else if tpe == "keep-alive" {
		return nil
	} else {
		return fmt.Errorf("Unknown message type: %s", tpe)
	}

	return nil
}

func (session *sessionStore) handleWelcome(conn *websocket.Conn, msg protocol.WelcomeRequest) error {
	session.conn[conn] = msg.Name
	conn.SetCloseHandler(session.closeHandler(conn))

	var pc protocol.TestParticipant
	if existingPc, ok := session.participants[msg.Name]; !ok {
        log.WithField("participant", msg.Name).Info("New participant joined")
		pc = protocol.TestParticipant{
			Name:         msg.Name,
			ClaimedCases: make(map[string]bool),
		}
		session.participants[msg.Name] = &pc
	} else {
        log.WithField("participant", msg.Name).Info("Old participant came back")
		pc = *existingPc
	}

	session.run.Participants = make([]protocol.TestParticipant, len(session.participants))
	idx := 0
	for _, nme := range session.participants {
		session.run.Participants[idx] = *nme
		idx++
	}
	session.Save()

	resp := protocol.WelcomeResponse{
		Type:        "welcome",
		Run:         *session.run,
		Suite:       *session.suite,
		Participant: pc,
	}
	return conn.WriteJSON(resp)
}

func (session *sessionStore) handleClaim(conn *websocket.Conn, msg protocol.ClaimRequest) error {
	name, ok := session.conn[conn]
	if !ok {
		return fmt.Errorf("No user associated with this connection. Did you say welcome?")
	}

	participant, ok := session.participants[name]
	if !ok {
		return fmt.Errorf("User %s seems not to participate in the testing. Looks like a bug.", name)
	}

	log.WithField("participant", participant.Name).WithField("case", msg.CaseID).WithField("claim", msg.Claim).Info("Participant claimed test case")
	if msg.Claim {
		participant.ClaimedCases[msg.CaseID] = true
	} else {
		delete(participant.ClaimedCases, msg.CaseID)
	}
	if err := session.Save(); err != nil {
		return err
	}

	resp := protocol.ClaimResponse{
		Type: "claim",
	}
	return conn.WriteJSON(resp)

	return nil
}

func (session *sessionStore) handleNewTestcaseRun(conn *websocket.Conn, msg protocol.NewTestCaseRunRequest) error {
	name, ok := session.conn[conn]
	if !ok {
		return fmt.Errorf("No user associated with this connection. Did you say welcome?")
	}

	participant, ok := session.participants[name]
	if !ok {
		return fmt.Errorf("User %s seems not to participate in the testing. Looks like a bug.", name)
	}

    tcr := protocol.TestCaseRun{
        Case: msg.Case,
        CaseGroup: msg.CaseGroup,
        Comment: msg.Comment,
        Result: msg.Result,
        Start: msg.Start,
        Tester: participant.Name,
    }
    log.WithField("participant", participant.Name).WithField("case", msg.Case).WithField("result", msg.Result).Info("Participant submitted testcase run")

    session.run.Cases = append(session.run.Cases, tcr)

	if err := session.Save(); err != nil {
		return err
	}

	resp := protocol.NewTestCaseRunResponse{
		Type: "newTestCaseRun",
	}
	return conn.WriteJSON(resp)

	return nil
}

func (session *sessionStore) Exit(conn *websocket.Conn) error {
	return nil
}

func (session *sessionStore) Save() error {
	fc, err := yaml.Marshal(session.run)
	if err != nil {
		return err
	}

	log.WithField("name", session.Name).Println("Wrote session log")
	if err := ioutil.WriteFile(session.Name, fc, 0644); err != nil {
		return err
	}

	for conn, name := range session.conn {
		participant := session.participants[name]
		err := conn.WriteJSON(protocol.UpdateMessage{
			Type:        "update",
			Participant: *participant,
			Run:         *(session.run),
		})
		if err != nil {
			log.WithError(err).Warn("Error while updating participant")
		}
	}
	return nil
}

func (session *sessionStore) closeHandler(conn *websocket.Conn) func(code int, text string) error {
	return func(code int, text string) error {
		pn, ok := session.conn[conn]
		if !ok {
			log.Warn("Connection without participant closed")
		}

		delete(session.conn, conn)
		log.WithField("participant", pn).Info("Participant left")
		return nil
	}
}
