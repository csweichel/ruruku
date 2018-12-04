package notifier

import (
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"sync"
)

type Notifier struct {
	mux  sync.Mutex
	cond *sync.Cond
	data *api.TestRunStatus
}

func NewNotifier() *Notifier {
	r := Notifier{}
	r.cond = sync.NewCond(&r.mux)
	return &r
}

func (n *Notifier) Listen() api.TestRunStatus {
	n.mux.Lock()
	n.cond.Wait()
	r := *n.data
	n.mux.Unlock()

	return r
}

func (n *Notifier) Update(data *api.TestRunStatus) {
	n.mux.Lock()
	n.data = data
	n.mux.Unlock()
	n.cond.Broadcast()
}
