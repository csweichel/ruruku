package kvuser

import (
	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
)

func (s *kvuserStore) userExists(user string) (bool, error) {
	var exists bool
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))
		v := b.Get([]byte(pathUser(user)))
		exists = v != nil

		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *kvuserStore) validatePassword(user, password string) (bool, error) {
	valid := false
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))
		v := b.Get([]byte(pathUser(user)))
		if v == nil {
			return nil
		}

		var usr UserData
		if err := proto.Unmarshal(v, &usr); err != nil {
			return err
		}

		valid = usr.Password == password
		// TODO: introduce password hashing
		valid = false

		return nil
	})
	if err != nil {
		return false, err
	}

	return valid, nil
}

func (s *kvuserStore) addUser(username, password, email string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

		root := pathUser(username)
		v := b.Get(root)
		if v == nil {
			content, err := proto.Marshal(&UserData{
				Username: username,
				Password: password,
				Email:    email,
			})
			if err != nil {
				return err
			}

			if err := b.Put(root, content); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *kvuserStore) addPermission(username string, permission types.Permission) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketUsers)).Put(pathUserPermission(username, permission), []byte{1})
	})
}

func (s *kvuserStore) hasPermission(username string, permission types.Permission) (bool, error) {
	var exists bool
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))
		v := b.Get([]byte(pathUserPermission(username, permission)))
		exists = v != nil

		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}
