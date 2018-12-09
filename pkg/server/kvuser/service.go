package kvuser

import (
	"context"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	"time"
)

const (
	bucketUsers = "users"
)

type kvuserStore struct {
	db            *bolt.DB
	tokenLifetime time.Duration
}

func NewUserStore(db *bolt.DB) (*kvuserStore, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketUsers))
		if err != nil {
			return err
		}

		root := pathUser("root")
		v := bucket.Get(root)
		if v == nil {
			content, err := proto.Marshal(&UserData{
				Username: "root",
				Password: "",
				Email:    "",
			})
			if err != nil {
				return err
			}

			if err := bucket.Put(root, content); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &kvuserStore{db: db, tokenLifetime: 30 * time.Minute}, nil
}

// AuthenticateCredentials authenticates a user based on username/password
func (s *kvuserStore) AuthenticateCredentials(ctx context.Context, req *api.AuthenticationRequest) (*api.AuthenticationRespose, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Add creates a new user with a set of credentials.
func (s *kvuserStore) Add(ctx context.Context, req *api.AddUserRequest) (*api.AddUserResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Delete removes an existing user. This invalidates all tokens of the user.
func (s *kvuserStore) Delete(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Grant adds permissions to a user
func (s *kvuserStore) Grant(ctx context.Context, req *api.GrantPermissionsRequest) (*api.GrantPermissionsResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// ChangePassword modifies the password of a user. This invalidates all tokens of the user.
func (s *kvuserStore) ChangePassword(ctx context.Context, req *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}
