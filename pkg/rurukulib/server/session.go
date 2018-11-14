package server

import (
    "fmt"
	"encoding/json"
	"github.com/32leaves/ruruku/protocol"
	"github.com/gorilla/websocket"
)

type sessionStore struct {
	Suite *protocol.TestSuite
	Run   *protocol.TestRun
    Name string
	conn  map[*websocket.Conn]string
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
    // TODO: handle name clash
    session.Run.Participants = append(session.Run.Participants, msg.Name)

    resp := protocol.WelcomeResponse {
        Type: "welcome",
        Run: *session.Run,
        Suite: *session.Suite,
    }
    return conn.WriteJSON(resp)
}

func (session *sessionStore) Exit(conn *websocket.Conn) (error) {
	return nil
}
