package progress

import (
	"fmt"
	"math"
	"time"

	"github.com/rockwell-uk/csync/mutex"
	"github.com/rockwell-uk/go-logger/logger"
	"github.com/rockwell-uk/go-utils/timeutils"
	"github.com/rockwell-uk/uiprogress"
)

type ProgressJob interface {
	Setup(jobName string, input interface{}) (*Job, error)
	Run(job *Job, input interface{}) (interface{}, error)
}

func RunJob(jobName, funcName string, job ProgressJob, magnitude int, setupInput interface{}, runInput interface{}) error {
	var start time.Time = time.Now()
	var took time.Duration

	logger.Log(
		logger.LVL_APP,
		fmt.Sprintf("%v [%v]\n", jobName, magnitude),
	)

	j, err := job.Setup(jobName, setupInput)
	if err != nil {
		return fmt.Errorf("%v %v", funcName, err.Error())
	}
	defer j.End(true)

	_, err = job.Run(j, runInput)
	if err != nil {
		return fmt.Errorf("%v %v", funcName, err.Error())
	}

	took = timeutils.Took(start)
	logger.Log(
		logger.LVL_DEBUG,
		fmt.Sprintf("Done %v [%v]\n", jobName, took),
	)

	return nil
}

type Job struct {
	Name      string
	Tasks     map[string]*Task
	StartTime *time.Time
	EndTime   *time.Time
	Took      *time.Duration
	Magnitude float64
	Progress  int // 0-100
	Bar       *uiprogress.Bar
}

func (j *Job) String() string {
	var s string

	s = fmt.Sprintf("Job: [%v]\n", j.Name)
	s += fmt.Sprintf("Magnitude: [%v]\n", j.Magnitude)
	s += fmt.Sprintf("Start: [%v]\n", timeutils.FormatTime(*j.GetStartTime()))

	for taskName, task := range j.GetTasks() {
		s += fmt.Sprintf("\nTask: [%v] %+v", taskName, *task)
	}

	return s
}

func (j *Job) AddTasks(tasks []*Task) {
	for _, task := range tasks {
		j.AddTask(task)
	}
}

func (j *Job) AddTask(t *Task) {
	j.GetTasks()[t.ID] = t
	j.CalculateMagnitude()
}

func (j *Job) GetTask(id string) (*Task, error) {
	if j == nil {
		return &Task{}, fmt.Errorf("task %s does not exist", id)
	}

	tasks := j.GetTasks()

	mutex.Lock()
	defer mutex.Unlock()
	if task, exists := tasks[id]; exists {
		return task, nil
	}

	return &Task{}, fmt.Errorf("task %s does not exist", id)
}

func (j *Job) CalculateMagnitude() {
	var magnitude float64 = 0

	for _, t := range j.GetTasks() {
		magnitude += t.Magnitude
	}
	if magnitude == 0 {
		panic(fmt.Sprintf("invalid magnitude for job %v", j.Name))
	}

	j.Magnitude = magnitude
}

func (j *Job) GetProgress() (float64, error) {
	t := j.GetStartTime()

	if t == nil {
		return 0, fmt.Errorf("job [%v] not started", j.Name)
	}

	done := j.GetCompletedMagnitude()
	if done == 0 {
		return 0, nil
	}

	var progress float64 = (done / j.Magnitude) * 100

	return math.Ceil(progress*100) / 100, nil
}

func (j *Job) GetDone() (int, error) {
	t := j.GetStartTime()

	if t == nil {
		return 0, fmt.Errorf("job [%v] not started", j.Name)
	}

	done := j.GetCompletedMagnitude()
	if done == 0 {
		return 0, nil
	}

	return int(done), nil
}

func (j *Job) GetRemaining() (string, error) {
	t := j.GetStartTime()

	if t == nil {
		return "100", fmt.Errorf("job [%v] not started", j.Name)
	}

	done := j.GetCompletedMagnitude()
	if done == 0 {
		return "???", nil
	}

	var remaining string

	elapsed := time.Since(*t)
	timePerUnit := float64(elapsed) / done
	remainingUnits := j.Magnitude - done
	remainingTime := time.Duration(int64(remainingUnits * timePerUnit))

	remaining = timeutils.Round(remainingTime, 0).String()

	return remaining, nil
}

func (j *Job) GetStatus() (string, error) {
	remaining, err := j.GetRemaining()
	if err != nil {
		return "", err
	}

	s := j.GetStartTime()
	t := j.GetTook()

	var elapsed string
	if t != nil {
		elapsed = timeutils.Round(*t, 0).String()
	} else {
		elapsed = timeutils.Round(time.Since(*s), 0).String()
	}

	return fmt.Sprintf("%v:%v", elapsed, remaining), nil
}

func (j *Job) GetCompletedMagnitude() float64 {
	var done float64

	for _, t := range j.GetTasks() {
		if t.GetTook() != nil {
			done += t.Magnitude
		}
	}

	return done
}

func (j *Job) GetStartTime() *time.Time {
	mutex.Lock()
	t := j.StartTime
	mutex.Unlock()

	return t
}

func (j *Job) GetTasks() map[string]*Task {
	mutex.Lock()
	t := j.Tasks
	mutex.Unlock()

	return t
}

func (j *Job) GetTook() *time.Duration {
	mutex.Lock()
	t := j.Took
	mutex.Unlock()

	return t
}

func (j *Job) Start() error {
	if len(j.GetTasks()) == 0 {
		return fmt.Errorf("tasks must be added to job [%v] before it can be started", j.Name)
	}

	t := time.Now()

	mutex.Lock()
	j.StartTime = &t
	mutex.Unlock()

	return nil
}

func (j *Job) End(stopProgress bool) {
	mutex.Lock()
	end := time.Now()
	j.EndTime = &end
	d := time.Since(*j.StartTime)
	j.Took = &d
	j.EndBar()
	mutex.Unlock()

	if stopProgress && ShouldShowBar() {
		uiprogress.Stop()
	}
}
