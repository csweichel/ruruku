package kvsession

import (
	"bytes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

const (
	bucketSessions = "Sessions"
	bucketTestplan = "Testplan"
)

func (s *kvsessionStore) isSession(sessionID string, getter func(s *SessionMetadata) bool) (bool, error) {
	var res bool
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))
		v := b.Get(pathSession(sessionID))
		if v == nil {
			return status.Errorf(codes.NotFound, "Session %s does not exist", sessionID)
		}

		var meta SessionMetadata
		if err := proto.Unmarshal(v, &meta); err != nil {
			return err
		}
		res = getter(&meta)

		return nil
	})
	return res, err
}

func (s *kvsessionStore) isSessionOpen(sessionID string) (bool, error) {
	return s.isSession(sessionID, func(session *SessionMetadata) bool {
		return session.Open
	})
}

func (s *kvsessionStore) isSessionModifiable(sessionID string) (bool, error) {
	return s.isSession(sessionID, func(session *SessionMetadata) bool {
		return session.Modifiable
	})
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

func (s *kvsessionStore) storeSession(sessionID string, name string, planID string, modifiable bool, annotations map[string]string) error {
	content, err := proto.Marshal(&SessionMetadata{
		Name:        name,
		PlanID:      planID,
		Open:        true,
		Modifiable:  modifiable,
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

func (s *kvsessionStore) modifySession(sessionID string, mod func(session *SessionMetadata)) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketSessions))

		v := b.Get(pathSession(sessionID))
		if v == nil {
			return status.Errorf(codes.NotFound, "Session %s does not exist", sessionID)
		}

		var meta SessionMetadata
		if err := proto.Unmarshal(v, &meta); err != nil {
			return err
		}

		mod(&meta)
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
				return status.Errorf(codes.Internal, "Cannot load session %s: %v", k, err)
			}

			err := cb(&api.ListSessionsResponse{
				Id:     getLastSegment(k),
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
