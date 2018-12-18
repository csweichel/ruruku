package kvsession

import (
	"fmt"

	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
)

func (s *kvsessionStore) claimTestcase(sessionID string, testcaseID string, userID string, claim bool) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketTestplan))
		if err != nil {
			return err
		}

		tckey := pathSessionTestcase(sessionID, testcaseID)
		if bucket.Get(tckey) == nil {
			return fmt.Errorf("Testcase '%s' does not exist", testcaseID)
		}

		key := pathSessionClaim(sessionID, testcaseID, userID)
		if claim {
			if err := bucket.Put([]byte(key), []byte{1}); err != nil {
				return err
			}
		} else {
			if err := bucket.Delete([]byte(key)); err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (s *kvsessionStore) hasClaimedTestcase(sessionID string, testcaseID string, userID string) (bool, error) {
	key := pathSessionClaim(sessionID, testcaseID, userID)

	var exists bool
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketTestplan))
		v := b.Get([]byte(key))
		exists = v != nil

		return nil
	})
	return exists, err
}

type kvsessionContribution struct {
	UserID     string
	TestcaseID string
	Result     types.TestRunState
	Comment    string
}

func (s *kvsessionStore) contribute(sessionID string, contribution kvsessionContribution) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketTestplan))
		if err != nil {
			return err
		}

		key := pathSessionContribution(sessionID, contribution.TestcaseID, contribution.UserID)
		content, err := proto.Marshal(&TestcaseContribution{
			Result:  string(contribution.Result),
			Comment: contribution.Comment,
		})
		if err != nil {
			return err
		}
		if err := bucket.Put(key, content); err != nil {
			return err
		}

		return nil
	})
	return err
}
