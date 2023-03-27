package progress

import (
	"fmt"
	"time"

	"github.com/rockwell-uk/csync/mutex"
)

type ProgressTask interface {
	Start()
	End()
}

type Task struct {
	ID        string
	Magnitude float64
	StartTime *time.Time
	Took      *time.Duration
}

func (t *Task) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t *Task) Start() {
	mutex.Lock()
	n := time.Now()
	t.StartTime = &n
	mutex.Unlock()
}

func (t *Task) End() {
	mutex.Lock()
	d := time.Since(*t.StartTime)
	t.Took = &d
	mutex.Unlock()
}

func (t *Task) GetTook() *time.Duration {
	mutex.Lock()
	k := t.Took
	mutex.Unlock()

	return k
}
