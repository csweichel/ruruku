package server

import (
	"sync"
    "os"
    "io/ioutil"
    "time"
    "github.com/32leaves/ruruku/protocol"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v2"
    "fmt"
)

type Storage interface {
    ClaimTestcase(tester string, caseID string, claim bool) error
    SetTestcaseRun(tester string, caseID string, result protocol.TestCaseResult, comment string) error
    AddParticipant(tester string) error
    Save() error

    GetSuite() protocol.TestSuite
    GetRun() protocol.TestRun
    GetParticipant(tester string) (protocol.TestParticipant, bool)
}

type OnSaveFunc func(st *FileStorage)

type FileStorage struct {
    OnSave OnSaveFunc

    fn string
    mux sync.Mutex
    suite        *protocol.TestSuite
	run          *protocol.TestRun

	runidx map[string]*protocol.TestCaseRun
    pptidx map[string]*protocol.TestParticipant
}

func LoadFileStorage(fn string, suite *protocol.TestSuite) (*FileStorage, error) {
    r := &FileStorage{
        fn: fn,
        suite: suite,
		pptidx: make(map[string]*protocol.TestParticipant),
		runidx: make(map[string]*protocol.TestCaseRun),
	}

	if _, err := os.Stat(fn); err != nil {
        return nil, err
    }

    fc, err := ioutil.ReadFile(fn)
    if err != nil {
        log.WithField("filename", fn).WithError(err).Error("Cannot read storage file")
        return nil, err
    }

    if err := yaml.Unmarshal(fc, r.run); err != nil {
        log.WithField("filename", fn).WithError(err).Error("Error while restoring storage")
        return nil, err
    }

    // build participant index
    for _, pcp := range r.run.Participants {
        r.pptidx[pcp.Name] = &pcp
    }

    // build run index
    for _, run := range r.run.Cases {
        r.runidx[fmt.Sprintf("%s#%s", run.Tester, run.CaseID)] = &run
    }

    return r, nil
}

func NewFileStorage(fn string, suite *protocol.TestSuite, runName string) *FileStorage {
    return &FileStorage{
        fn: fn,
		suite: suite,
		run: &protocol.TestRun{
			Name:      runName,
			SuiteName: suite.Name,
			Start:     time.Now().Format(time.RFC3339),
		},
		pptidx: make(map[string]*protocol.TestParticipant),
		runidx:       make(map[string]*protocol.TestCaseRun),
	}
}

func (st *FileStorage) ClaimTestcase(tester string, caseID string, claim bool) error {
    st.mux.Lock()
    defer st.mux.Unlock()

    participant, ok := st.pptidx[tester]
	if !ok {
		return fmt.Errorf("User %s seems not to participate in the testing. Looks like a bug.", tester)
	}

	log.WithField("participant", participant.Name).WithField("case", caseID).WithField("claim", claim).Info("Participant claimed test case")
	if claim {
		participant.ClaimedCases[caseID] = true
	} else {
		delete(participant.ClaimedCases, caseID)
	}

    return nil
}

func (st *FileStorage) SetTestcaseRun(tester string, caseID string, result protocol.TestCaseResult, comment string) error {
    st.mux.Lock()
    defer st.mux.Unlock()

    participant, ok := st.pptidx[tester]
	if !ok {
		return fmt.Errorf("User %s seems not to participate in the testing. Looks like a bug.", tester)
	}

	// check if we have this run in the index already
	idxkey := fmt.Sprintf("%s#%s", participant.Name, caseID)
	if run, ok := st.runidx[idxkey]; ok {
		run.Comment = comment
		run.Result = result
	} else {
		tcr := protocol.TestCaseRun{
			CaseID:  caseID,
			Comment: comment,
			Result:  result,
			Start:   time.Now().Format(time.RFC3339),
			Tester:  participant.Name,
		}
		st.runidx[idxkey] = &tcr
	}

	st.run.Cases = make([]protocol.TestCaseRun, len(st.runidx))
	i := 0
	for _, run := range st.runidx {
		st.run.Cases[i] = *run
		i++
	}

	log.WithField("participant", participant.Name).WithField("case", caseID).WithField("result", result).Info("Participant submitted testcase run")

    return nil
}

func (st *FileStorage) AddParticipant(tester string) error {
    st.mux.Lock()
    defer st.mux.Unlock()

    if _, ok := st.pptidx[tester]; ok {
        return fmt.Errorf("Participant %s is already registered", tester)
    }

    pc := protocol.TestParticipant{
        Name:         tester,
        ClaimedCases: make(map[string]bool),
    }
    st.pptidx[tester] = &pc

    st.run.Participants = make([]protocol.TestParticipant, len(st.pptidx))
	idx := 0
	for _, nme := range st.pptidx {
		st.run.Participants[idx] = *nme
		idx++
	}

    return nil
}

func (st *FileStorage) Save() error {
    st.mux.Lock()
    defer st.mux.Unlock()

    fc, err := yaml.Marshal(st.run)
	if err != nil {
		return err
	}

	log.WithField("name", st.fn).Println("Wrote session log")
	if err := ioutil.WriteFile(st.fn, fc, 0644); err != nil {
		return err
	}

    if st.OnSave != nil {
        st.OnSave(st)
    }

    return nil
}

func (st *FileStorage) GetSuite() protocol.TestSuite {
    return *st.suite
}

func (st *FileStorage) GetRun() protocol.TestRun {
    return *st.run
}

func (st *FileStorage) GetParticipant(tester string) (protocol.TestParticipant, bool) {
    if r, ok := st.pptidx[tester]; ok {
        return *r, true
    } else {
        return protocol.TestParticipant{}, false
    }
}
