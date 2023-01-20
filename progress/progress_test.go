package progress

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestJob(t *testing.T) {

	n := 10

	j := NewJob("TestCalcProgress", n)

	// Add n tasks with arbitrary magnitude
	for i := 0; i < n; i++ {
		m := 100.0 * float64(i+1)
		j.AddTask(&Task{
			ID:        fmt.Sprintf("task_%v", i),
			Magnitude: m,
		})
	}

	err := j.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Fudge the start time of the job
	s := j.StartTime.Add(-time.Second * 50)
	j.StartTime = &s

	expected := []string{
		"0% complete, remaining ???",
		"1.82% complete, remaining 45m0s",
		"5.46% complete, remaining 14m27s",
		"10.91% complete, remaining 6m48s",
		"18.19% complete, remaining 3m45s",
		"27.28% complete, remaining 2m13s",
		"38.19% complete, remaining 1m21s",
		"50.91% complete, remaining 48s",
		"65.46% complete, remaining 26s",
		"81.82% complete, remaining 11s",
		"100% complete, remaining 0s",
	}
	actual := []string{}

	progress, err := j.GetProgress()
	if err != nil {
		t.Fatal(err)
	}

	remaining, err := j.GetRemaining()
	if err != nil {
		t.Fatal(err)
	}

	actual = append(actual, fmt.Sprintf("%v%s complete, remaining %v", progress, "%", remaining))

	for i := 0; i < n; i++ {

		task, _ := j.GetTask(fmt.Sprintf("task_%v", i))

		task.Start()

		// Fake task.End()
		s := i * 100
		took := (time.Second * time.Duration(s))
		task.Took = &took

		progress, err := j.GetProgress()
		if err != nil {
			t.Fatal(err)
		}

		remaining, err := j.GetRemaining()
		if err != nil {
			t.Fatal(err)
		}

		actual = append(actual, fmt.Sprintf("%v%s complete, remaining %v", progress, "%", remaining))
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, actual %v", prettyPrint(expected), prettyPrint(actual))
	}
}

func TestGetProgress(t *testing.T) {

	n := 10
	m := 100.0
	took := time.Second * 10

	j := NewJob("TestGetProgress", n)

	// Setup an in progress job
	for i := 0; i < n; i++ {
		if i < 5 {
			j.AddTask(&Task{
				ID:        fmt.Sprintf("task_%v", i),
				Magnitude: m,
				Took:      &took,
			})
		} else {
			j.AddTask(&Task{
				ID:        fmt.Sprintf("task_%v", i),
				Magnitude: m,
			})
		}
	}

	err := j.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Fudge the start time of the job
	s := j.StartTime.Add(-time.Second * 50)
	j.StartTime = &s

	expected := "50% complete, remaining 50s"

	progress, err := j.GetProgress()
	if err != nil {
		t.Fatal(err)
	}

	remaining, err := j.GetRemaining()
	if err != nil {
		t.Fatal(err)
	}

	actual := fmt.Sprintf("%v%s complete, remaining %v", progress, "%", remaining)

	if expected != actual {
		t.Errorf("expected [%#v], actual [%#v]", expected, actual)
	}
}

func TestStartJobError(t *testing.T) {

	jobName := "TestStartJobError"
	n := 10
	j := NewJob(jobName, n)

	err := j.Start()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := fmt.Errorf("tasks must be added to job [%v] before it can be started", jobName)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, actual %v", expected, err)
	}
}

func TestGetProgressError(t *testing.T) {

	jobName := "TestGetProgressError"
	n := 10
	j := NewJob(jobName, n)

	// Add n tasks with arbitrary magnitude
	for i := 1; i <= n; i++ {
		m := 100.0 * float64(i)
		j.AddTask(&Task{
			ID:        fmt.Sprintf("task_%v", i),
			Magnitude: m,
		})
	}

	_, err := j.GetProgress()
	expected := fmt.Errorf("job [%v] not started", jobName)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, actual %v", expected, err)
	}
}

func TestGetRemainingError(t *testing.T) {

	jobName := "TestGetProgressError"
	n := 10
	j := NewJob(jobName, n)

	// Add n tasks with arbitrary magnitude
	for i := 1; i <= n; i++ {
		m := 100.0 * float64(i)
		j.AddTask(&Task{
			ID:        fmt.Sprintf("task_%v", i),
			Magnitude: m,
		})
	}

	_, err := j.GetRemaining()
	expected := fmt.Errorf("job [%v] not started", jobName)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, actual %v", expected, err)
	}
}

func prettyPrint(data interface{}) string {

	var p []byte

	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(p)
}
