package kvsession

import (
	"bytes"
	"fmt"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

const (
	bucketSessions = "Sessions"
	bucketTestplan = "Testplan"
)

func (s *kvsessionStore) isSessionOpen(sessionID string) (bool, error) {
	var open bool
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))
		v := b.Get(pathSession(sessionID))
		if v == nil {
			return fmt.Errorf("Session %s does not exist", sessionID)
		}

		var meta SessionMetadata
		if err := proto.Unmarshal(v, &meta); err != nil {
			return err
		}
		open = meta.Open

		return nil
	})
	return open, err
}

func (s *kvsessionStore) sessionExists(sessionID string) (bool, error) {
	var exists bool
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))
		v := b.Get(pathSession(sessionID))
		exists = v != nil

		return nil
	})
	return exists, err
}

func (s *kvsessionStore) storeSession(sessionID string, name string, planID string, annotations map[string]string) error {
	content, err := proto.Marshal(&SessionMetadata{
		Name:        name,
		PlanID:      planID,
		Open:        true,
		Annotations: annotations,
	})
	if err != nil {
		return err
	}

	err = s.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketSessions))
		if err != nil {
			return err
		}

		return bucket.Put(pathSession(sessionID), content)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *kvsessionStore) closeSession(sessionID string) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))

		v := b.Get(pathSession(sessionID))
		if v == nil {
			return fmt.Errorf("Session %s does not exist", sessionID)
		}

		var meta SessionMetadata
		if err := proto.Unmarshal(v, &meta); err != nil {
			return err
		}

		meta.Open = false
		content, err := proto.Marshal(&meta)
		if err != nil {
			return err
		}

		return b.Put(pathSession(sessionID), content)
	})
}

func (s *kvsessionStore) listSessions(cb func(session *api.ListSessionsResponse) error) error {
	return s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))
		c := b.Cursor()

		prefix := pathSessions()
		var r SessionMetadata
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			if err := proto.Unmarshal(v, &r); err != nil {
				return fmt.Errorf("Cannot load session %s: %v", k, err)
			}

			err := cb(&api.ListSessionsResponse{
				Id:     string(k),
				Name:   r.Name,
				IsOpen: r.Open,
			})
			if err != nil {
				log.WithError(err).Debug("Error from listSession callback")
				return err
			}
		}

		return nil
	})
}

func (s *kvsessionStore) registerParticipant(sessionID string, name string) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))

		key := pathSessionParticipant(sessionID, name)
		return b.Put(key, []byte(name))
	})
}

func (s *kvsessionStore) getParticipant(sessionID string, participantID string) (*api.Participant, error) {
	key := pathSessionParticipant(sessionID, participantID)

	var result *api.Participant = nil
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))
		v := b.Get([]byte(key))

		if v != nil {
			result = &api.Participant{Name: string(v)}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *kvsessionStore) participantInSession(sessionID string, participantID string) (bool, error) {
	key := pathSessionParticipant(sessionID, participantID)

	var exists bool
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))
		v := b.Get([]byte(key))
		exists = v != nil

		return nil
	})

	return exists, err
}
