package server

import (
	"encoding/json"
	"fmt"
	"github.com/32leaves/ruruku/protocol"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type sessionStore struct {
	suite        *protocol.TestSuite
	run          *protocol.TestRun
	Name         string
	conn         map[*websocket.Conn]string
	participants map[string]protocol.TestParticipant
}

func NewSession(name string, suite *protocol.TestSuite) *sessionStore {
	return &sessionStore{
		suite: suite,
		run: &protocol.TestRun{
			SuiteName: suite.Name,
			Start:     time.Now().Format(time.RFC3339),
		},
		conn:         make(map[*websocket.Conn]string),
		Name:         name,
		participants: make(map[string]protocol.TestParticipant),
	}
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
	} else {
		return fmt.Errorf("Unknown message type: %s", tpe)
	}

	return nil
}

func (session *sessionStore) handleWelcome(conn *websocket.Conn, msg protocol.WelcomeRequest) error {
	session.conn[conn] = msg.Name

	var pc protocol.TestParticipant
	if existingPc, ok := session.participants[msg.Name]; !ok {
		pc = protocol.TestParticipant{
			Name:         msg.Name,
			ClaimedCases: make([]string, 0),
		}
		session.participants[msg.Name] = pc
	} else {
		pc = existingPc
	}

	session.run.Participants = make([]protocol.TestParticipant, len(session.participants))
	idx := 0
	for _, nme := range session.participants {
		session.run.Participants[idx] = nme
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

func (session *sessionStore) Exit(conn *websocket.Conn) error {
	return nil
}

func (session *sessionStore) Save() error {
	fc, err := yaml.Marshal(session.run)
	if err != nil {
		return err
	}

	log.WithField("name", session.Name).Println("Wrote session log")
	return ioutil.WriteFile(session.Name, fc, 0644)
}
