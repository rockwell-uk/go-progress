package progress

import (
	"github.com/rockwell-uk/go-logger/logger"
	"github.com/rockwell-uk/uiprogress"
)

func SetupJob(jobName string, tasks []*Task) *Job {
	job := NewJob(jobName, len(tasks))
	job.AddTasks(tasks)
	job.CalculateMagnitude()

	err := job.Start()
	if err != nil {
		panic(err)
	}

	if ShouldShowBar() {
		uiprogress.Start()

		job.Bar = uiprogress.AddBar(int(job.Magnitude)).PrependCompleted()
		job.Bar.AppendFunc(func(bar *uiprogress.Bar) string {
			status, err := job.GetStatus()
			if err != nil {
				panic(err)
			}
			return status
		})
	}

	return job
}

func NewJob(name string, numTasks int) *Job {
	return &Job{
		Name:  name,
		Tasks: make(map[string]*Task, numTasks),
	}
}

func ShouldShowBar() bool {
	return logger.Vbs == logger.LVL_APP
}
