package kvuser

import (
	"fmt"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"time"
)

// GetUserToken returns a fresh token for a user
func (s *kvuserStore) GetUserToken(user string) (string, error) {
	if exists, err := s.userExists(user); err != nil {
		return "", err
	} else if !exists {
		return "", fmt.Errorf("User %s does not exist", user)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("Cannot create token: %v", err)
	}
	token := id.String()

	content, err := proto.Marshal(&TokenData{
		Username:  user,
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return "", err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketUsers)).Put(pathUserToken(token), content)
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *kvuserStore) getUserFromToken(token string) (string, error) {
	var username string
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

		key := pathUserToken(token)
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}

		var td TokenData
		if err := proto.Unmarshal(v, &td); err != nil {
			return err
		}

		age := time.Now().Sub(time.Unix(td.Timestamp, 0))
		if age > s.TokenLifetime {
			b.Delete(key)
			return nil
		}

		uk := pathUser(td.Username)
		if b.Get(uk) == nil {
			// user no longer exists
			b.Delete(key)
			return nil
		}

		username = td.Username
		return nil
	})
	if err != nil {
		return "", err
	}

	return username, nil
}
