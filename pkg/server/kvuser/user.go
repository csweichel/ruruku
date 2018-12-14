package kvuser

import (
	"bytes"
	"fmt"
	"github.com/32leaves/ruruku/pkg/types"
    api "github.com/32leaves/ruruku/pkg/api/v1"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/bcrypt"
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

		valid = bcrypt.CompareHashAndPassword(usr.Password, []byte(password)) == nil

		return nil
	})
	if err != nil {
		return false, err
	}

	return valid, nil
}

func (s *kvuserStore) changePassword(username, password string) error {
	pwdhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

		root := pathUser(username)
		v := b.Get(root)
		if v != nil {
			var usr UserData
			if err := proto.Unmarshal(v, &usr); err != nil {
				return err
			}

			content, err := proto.Marshal(&UserData{
				Username: username,
				Password: pwdhash,
				Email:    usr.Email,
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

func (s *kvuserStore) addUser(username, password, email string) error {
	pwdhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

		root := pathUser(username)
		v := b.Get(root)
		if v == nil {
			content, err := proto.Marshal(&UserData{
				Username: username,
				Password: pwdhash,
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

func (s *kvuserStore) addPermissions(username string, permission []types.Permission) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))
		for _, perm := range permission {
			if err := b.Put(pathUserPermission(username, perm), []byte(string(perm))); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *kvuserStore) deleteUser(username string) error {
	if username == "root" {
		return fmt.Errorf("cannot delete root")
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

		if err := b.Delete(pathUser(username)); err != nil {
			return err
		}
		prefix := pathUserPermissions(username)
		c := b.Cursor()
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			if err := b.Delete(k); err != nil {
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

func (s *kvuserStore) hasPermission(username string, permission types.Permission) (bool, error) {
	if username == "root" {
		return true, nil
	}

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

func (s *kvuserStore) listUsers() ([]*api.User, error) {
    result := make([]*api.User, 0)
    err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

        prefix := []byte("u")
        c := b.Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
            var usr UserData
			if err := proto.Unmarshal(v, &usr); err != nil {
				return err
			}

            permissions, err := s.listPermissions(usr.Username)
            if err != nil {
                return err
            }

            user := api.User{
                Name: usr.Username,
                Email: usr.Email,
                Permission: permissions,
            }
            result = append(result, &user)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

    return result, nil
}

func (s *kvuserStore) listPermissions(user string) ([]api.Permission, error) {
    if user == "root" {
        r := make([]api.Permission, len(types.AllPermissions))
        for idx, p := range types.AllPermissions {
            r[idx] = api.ConvertPermission(p)
        }
        return r, nil
    }

    result := make([]api.Permission, 0)
    err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketUsers))

        prefix := pathUserPermissions(user)
        c := b.Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
            result = append(result, api.ConvertPermission(types.Permission(v)))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

    return result, nil
}