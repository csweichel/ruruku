package kvsession

import (
	"bytes"
	"fmt"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *kvsessionStore) storePlan(sessionID string, plan *api.TestPlan) error {
	planmeta, err := proto.Marshal(&TestplanMetadata{
		Id:          plan.Id,
		Name:        plan.Name,
		Description: plan.Description,
	})
	if err != nil {
		return err
	}

	err = s.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketTestplan))
		if err != nil {
			return err
		}

		plankey := pathSessionPlan(sessionID)
		if err := bucket.Put([]byte(plankey), planmeta); err != nil {
			return err
		}

		for _, cse := range plan.Case {
			tckey := pathSessionTestcase(sessionID, cse.Id)
			caseExists := bucket.Get(tckey) != nil
			if caseExists {
				return status.Errorf(codes.AlreadyExists, "Testcase '%s' exists already", cse.Id)
			}

			content, err := proto.Marshal(cse)
			if err != nil {
				return err
			}
			if err := bucket.Put(tckey, content); err != nil {
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

func (s *kvsessionStore) addOrUpdateTestcase(sessionID string, tc []*api.Testcase, tcMustExist bool) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketTestplan))
		if err != nil {
			return err
		}

		for _, cse := range tc {
			tckey := pathSessionTestcase(sessionID, cse.Id)
			caseExists := bucket.Get(tckey) != nil
			if tcMustExist && !caseExists {
				return status.Errorf(codes.NotFound, "Testcase '%s' does not exist", cse.Id)
			}
			if !tcMustExist && caseExists {
				return status.Errorf(codes.AlreadyExists, "Testcase '%s' exists already", cse.Id)
			}

			content, err := proto.Marshal(cse)
			if err != nil {
				return err
			}
			if err := bucket.Put(tckey, content); err != nil {
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

func (s *kvsessionStore) removeTestcase(sessionID string, tc []*api.Testcase) error {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketTestplan))
		if err != nil {
			return err
		}

		for _, cse := range tc {
			tckey := pathSessionTestcase(sessionID, cse.Id)
			caseExists := bucket.Get(tckey) != nil
			if !caseExists {
				return status.Errorf(codes.NotFound, "Testcase '%s' does not exist", cse.Id)
			}

			if err := bucket.Delete(tckey); err != nil {
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

func (s *kvsessionStore) testcaseExists(sessionID string, testcaseID string) (bool, error) {
	var exists bool
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketTestplan))
		v := b.Get(pathSessionTestcase(sessionID, testcaseID))
		exists = v != nil

		return nil
	})
	return exists, err
}

func (s *kvsessionStore) getStatus(sessionID string) (*api.TestRunStatus, error) {
	res := api.TestRunStatus{}
	sessionState := types.Passed
	err := s.DB.View(func(tx *bolt.Tx) error {
		sv := tx.Bucket([]byte(bucketSessions)).Get(pathSession(sessionID))
		if sv == nil {
			return fmt.Errorf("Session %s does not exist", sessionID)
		}

		var meta SessionMetadata
		if err := proto.Unmarshal(sv, &meta); err != nil {
			return err
		}
		res.Id = sessionID
		res.PlanID = meta.PlanID
		res.Name = meta.Name
		res.Open = meta.Open
		res.Modifiable = meta.Modifiable
		if meta.Annotations == nil {
			res.Annotations = map[string]string{}
		} else {
			res.Annotations = meta.Annotations
		}

		cases := make([]*api.TestcaseStatus, 0)
		tb := tx.Bucket([]byte(bucketTestplan))
		err := forEachTestcase(tb, sessionID, func(tcid string, tc *api.Testcase) error {
			claims := make([]*api.Participant, 0)
			err := forEachClaim(tb, sessionID, tcid, func(uid string) error {
				p, err := s.getParticipant(sessionID, uid)
				if err != nil {
					return err
				}
				claims = append(claims, p)

				return nil
			})
			if err != nil {
				return err
			}

			tcres := types.Passed
			results := make([]*api.TestcaseRunResult, 0)
			err = forEachResult(tb, sessionID, tcid, func(uid string, contrib *TestcaseContribution) error {
				p, err := s.getParticipant(sessionID, uid)
				if err != nil {
					return err
				}

				res := types.TestRunState(contrib.Result)
				tcres = types.WorseState(tcres, res)

				results = append(results, &api.TestcaseRunResult{
					Participant: p,
					State:       api.ConvertTestRunState(res),
					Comment:     contrib.Comment,
				})
				return nil
			})
			if err != nil {
				return err
			}
			if len(results) == 0 {
				tcres = types.Undecided
			}

			sessionState = types.WorseState(sessionState, tcres)

			thistc := *tc
			cse := api.TestcaseStatus{
				Case:   &thistc,
				Claim:  claims,
				Result: results,
				State:  api.ConvertTestRunState(tcres),
			}

			cases = append(cases, &cse)
			return nil
		})
		if err != nil {
			return err
		}

		res.Case = cases
		return nil
	})
	if err != nil {
		return nil, err
	}

	res.State = api.ConvertTestRunState(sessionState)
	return &res, nil
}

func forEachTestcase(b *bolt.Bucket, sessionID string, cb func(tcid string, tc *api.Testcase) error) error {
	var tc api.Testcase
	c := b.Cursor()
	prefix := pathSessionTestcases(sessionID)
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		tcid := getLastSegment(k)
		if err := proto.Unmarshal(v, &tc); err != nil {
			return err
		}
		if err := cb(tcid, &tc); err != nil {
			return err
		}
	}
	return nil
}

func forEachClaim(b *bolt.Bucket, sessionID string, testcaseID string, cb func(userID string) error) error {
	c := b.Cursor()
	prefix := pathSessionClaims(sessionID, testcaseID)
	for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
		uid := getLastSegment(k)
		if err := cb(uid); err != nil {
			return err
		}
	}
	return nil
}

func forEachResult(b *bolt.Bucket, sessionID string, testcaseID string, cb func(userID string, res *TestcaseContribution) error) error {
	var r TestcaseContribution
	c := b.Cursor()
	prefix := pathSessionContributions(sessionID, testcaseID)
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		uid := getLastSegment(k)
		if err := proto.Unmarshal(v, &r); err != nil {
			return err
		}
		if err := cb(uid, &r); err != nil {
			return err
		}
	}
	return nil
}
