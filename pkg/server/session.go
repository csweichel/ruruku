package server

import (
	"encoding/json"
	"fmt"
	"github.com/32leaves/ruruku/protocol"
    "github.com/32leaves/ruruku/pkg/storage"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"os"
    "sync"
)

type sessionStore struct {
	conn         map[*websocket.Conn]string
    mux sync.Mutex
	st *storage.FileStorage
}

func LoadSessionOrNew(name string, suite *protocol.TestSuite) (*sessionStore, error) {
    var st *storage.FileStorage
    fn := fmt.Sprintf("%s.yaml", name)
    if _, err := os.Stat(fn); err == nil {
        st, err = storage.LoadFileStorage(fn, suite)
        if err != nil {
            return nil, err
        }
    } else {
        st = storage.NewFileStorage(fn, suite, name)
    }

	r := &sessionStore{
		conn:         make(map[*websocket.Conn]string),
        st: st,
	}
    st.OnSave = r.updateClients

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
    session.mux.Lock()
	session.conn[conn] = msg.Name
    session.mux.Unlock()

	conn.SetCloseHandler(session.closeHandler(conn))

    if err := session.st.AddParticipant(msg.Name); err == nil {
		log.WithField("participant", msg.Name).Info("New participant joined")
        if err := session.st.Save(); err != nil {
            log.WithError(err).Warn("Unable to persist session")
        }
    } else {
		log.WithField("participant", msg.Name).Info("Old participant came back")
    }

    p, _ := session.st.GetParticipant(msg.Name)
	resp := protocol.WelcomeResponse{
		Type:        "welcome",
		Run:         session.st.GetRun(),
		Suite:       session.st.GetSuite(),
        Participant: p,
	}
	return conn.WriteJSON(resp)
}

func (session *sessionStore) handleClaim(conn *websocket.Conn, msg protocol.ClaimRequest) error {
    session.mux.Lock()
    tester, ok := session.conn[conn]
    session.mux.Unlock()
	if !ok {
		return fmt.Errorf("No user associated with this connection. Did you say welcome?")
	}

    if err := session.st.ClaimTestcase(tester, msg.CaseID, msg.Claim); err != nil {
        return err
    }

	if err := session.st.Save(); err != nil {
		return err
	}

	resp := protocol.ClaimResponse{
		Type: "claim",
	}
	return conn.WriteJSON(resp)

	return nil
}

func (session *sessionStore) handleNewTestcaseRun(conn *websocket.Conn, msg protocol.NewTestCaseRunRequest) error {
    session.mux.Lock()
    tester, ok := session.conn[conn]
    session.mux.Unlock()
	if !ok {
		return fmt.Errorf("No user associated with this connection. Did you say welcome?")
	}

    if err := session.st.SetTestcaseRun(tester, msg.CaseID, msg.Result, msg.Comment); err != nil {
        return err
    }

	log.WithField("participant", tester).WithField("case", msg.CaseID).WithField("result", msg.Result).Info("Participant submitted testcase run")

	if err := session.st.Save(); err != nil {
		return err
	}

	resp := protocol.NewTestCaseRunResponse{
		Type: "newTestCaseRun",
	}
	if err := conn.WriteJSON(resp); err != nil {
        return err
    }

    return nil
}

func (session *sessionStore) Exit(conn *websocket.Conn) error {
    session.mux.Lock()
	if name, ok := session.conn[conn]; ok {
		log.WithField("participant", name).Info("Participant exiting")
	} else {
		log.Warn("Exiting session with unknown participant")
	}

	delete(session.conn, conn)
    session.mux.Unlock()

	return conn.Close()
}

func (session *sessionStore) updateClients(st *storage.FileStorage) {
    session.mux.Lock()
    run := session.st.GetRun()
    for conn, name := range session.conn {
		participant, _ := session.st.GetParticipant(name)
		err := conn.WriteJSON(protocol.UpdateMessage{
			Type:        "update",
			Participant: participant,
			Run:         run,
		})
		if err != nil {
			log.WithError(err).Warn("Error while updating participant - dropping participant")
            session.mux.Unlock()
			session.Exit(conn)
            session.mux.Lock()
		}
	}
    session.mux.Unlock()
}

func (session *sessionStore) closeHandler(conn *websocket.Conn) func(code int, text string) error {
	return func(code int, text string) error {
        return session.Exit(conn)
	}
}
