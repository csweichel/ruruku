package server

import (
	"fmt"
	"github.com/32leaves/ruruku/protocol"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var session *sessionStore

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // not checking origin
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./client/build"+r.URL.Path)
}

// this is also the handler for joining to the chat
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %s", err)
		return
	}

	go func() {
		log.Info("New websocket connection")

		if err := session.Join(conn); err != nil {
			log.WithError(err).Error("Unable to join session")
			return
		}

		// then watch for incoming messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil { // if error then assuming that the connection is closed
				log.WithError(err).Error("Error while reading message from WS")
				session.Exit(conn)
				return
			}

			if err := session.HandleMessage(conn, message); err != nil {
				log.WithError(err).Error("Error while handling message")
				return
			}
		}

	}()
}

func Start(cfg *Config, suite *protocol.TestSuite, sessionName string) error {
	if cfg.Token == "" {
		cfg.Token = uuid.Must(uuid.NewV4()).String()
	}

	session = NewSession(sessionName, suite)

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", staticFiles)

	// fmt.Println("\nSuccess! Please navigate your browser to http://localhost:8000")
	log.Printf("Server started: http://localhost:%d/?token=%s", cfg.Port, cfg.Token)
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
