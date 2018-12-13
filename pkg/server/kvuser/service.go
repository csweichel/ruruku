package kvuser

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
	"time"
)

const (
	bucketUsers = "users"
)

type kvuserStore struct {
	db            *bolt.DB
	TokenLifetime time.Duration
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
				Password: []byte{},
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

	return &kvuserStore{db: db, TokenLifetime: 30 * time.Minute}, nil
}

// AuthenticateCredentials authenticates a user based on username/password
func (s *kvuserStore) AuthenticateCredentials(ctx context.Context, req *api.AuthenticationRequest) (*api.AuthenticationRespose, error) {
	if req.Username == "root" {
		return nil, status.Error(codes.Unauthenticated, "cannot authenticate as root")
	}

	if exists, err := s.userExists(req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !exists {
		return nil, status.Errorf(codes.NotFound, "unknown user %s", req.Username)
	}

	if valid, err := s.validatePassword(req.Username, req.Password); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "cannot authenticate %s", req.Username)
	}

	tkn, err := s.GetUserToken(req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.WithField("user", req.Username).Debug("Successfully authenticated user")

	return &api.AuthenticationRespose{Token: tkn}, nil
}

// Add creates a new user with a set of credentials.
func (s *kvuserStore) Add(ctx context.Context, req *api.AddUserRequest) (*api.AddUserResponse, error) {
	if err := s.ValidUserFromRequest(ctx, types.PermissionUserAdd); err != nil {
		return nil, err
	}

	if exists, err := s.userExists(req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if exists {
		return nil, status.Errorf(codes.AlreadyExists, "user %s exists already", req.Username)
	}

	if _, err := mail.ParseAddress(req.Email); req.Username == "" || len(req.Password) < 4 || err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user data is invalid: %v", err)
	}

	if err := s.addUser(req.Username, req.Password, req.Email); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.AddUserResponse{}, nil
}

// Delete removes an existing user. This invalidates all tokens of the user.
func (s *kvuserStore) Delete(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	if err := s.ValidUserFromRequest(ctx, types.PermissionUserDelete); err != nil {
		return nil, err
	}

	if req.Username == "root" {
		return nil, status.Error(codes.PermissionDenied, "Cannot delete root")
	}

	if exists, err := s.userExists(req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !exists {
		return nil, status.Errorf(codes.NotFound, "user %s does not exist", req.Username)
	}

	if err := s.deleteUser(req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.DeleteUserResponse{}, nil
}

// Grant adds permissions to a user
func (s *kvuserStore) Grant(ctx context.Context, req *api.GrantPermissionsRequest) (*api.GrantPermissionsResponse, error) {
	if err := s.ValidUserFromRequest(ctx, types.PermissionUserGrant); err != nil {
		return nil, err
	}

	if req.Username == "root" {
		return nil, status.Error(codes.PermissionDenied, "Cannot grant permissions to root")
	}

	if exists, err := s.userExists(req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !exists {
		return nil, status.Errorf(codes.NotFound, "user %s does not exist", req.Username)
	}

	permissions := make([]types.Permission, len(req.Permission))
	for idx, p := range req.Permission {
		permissions[idx] = p.Convert()
	}
	if err := s.addPermissions(req.Username, permissions); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.WithField("user", req.Username).WithField("permissions", req.Permission).Debug("Granted permissions")

	return &api.GrantPermissionsResponse{}, nil
}

// ChangePassword modifies the password of a user. This invalidates all tokens of the user.
func (s *kvuserStore) ChangePassword(ctx context.Context, req *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	usr, err := s.getUserFromRequest(ctx)
	if err != nil {
		return nil, err
	}

	if ok, err := s.hasPermission(usr, types.PermissionUserChpwd); (err != nil || !ok) && usr != req.Username {
		return nil, status.Errorf(codes.PermissionDenied, "User does not have %v permission", types.PermissionUserChpwd)
	}

	if req.Username == "root" {
		return nil, status.Error(codes.PermissionDenied, "Cannot change password of root")
	}

	if exists, err := s.userExists(req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else if !exists {
		return nil, status.Errorf(codes.NotFound, "user %s does not exist", req.Username)
	}

	if err := s.changePassword(req.Username, req.NewPassword); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.ChangePasswordResponse{}, nil
}
