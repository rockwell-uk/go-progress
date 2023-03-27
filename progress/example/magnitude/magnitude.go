//nolint
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/rockwell-uk/csync/counter"
	"github.com/rockwell-uk/csync/waitgroup"
	"github.com/rockwell-uk/go-logger/logger"
	"github.com/rockwell-uk/go-progress/progress"
	"github.com/rockwell-uk/uiprogress"
)

var (
	chunkSize float64 = 1000

	testFiles = []string{
		"test1.sql",
		"test2.sql",
		"test3.sql",
	}

	magnitudes = map[string]float64{
		"test1.sql": 100000,
		"test2.sql": 200000,
		"test3.sql": 150000,
	}

	ctr *counter.Counter
)

func main() {

	logger.Start(logger.LVL_APP)

	var jobName string = "testMagnitude"

	ctr = counter.New()

	uiprogress.Start()

	var totalChunks int
	var tasks []*progress.Task
	for _, size := range magnitudes {
		for i := 0; i <= int(math.Ceil(size/chunkSize)); i++ {
			tasks = append(tasks, &progress.Task{
				ID:        chunkName(totalChunks),
				Magnitude: 1,
			})
			totalChunks++
		}
	}

	job := progress.SetupJob(jobName, tasks)
	job.Start()

	fanout(testFiles, job)

	job.End(true)

	fmt.Printf("%v ended\n", jobName)
}

func fanout(testFiles []string, job *progress.Job) {

	var wg *waitgroup.WaitGroup = waitgroup.New()

	for _, testFile := range testFiles {
		wg.Add(1)
		go func(j *progress.Job, sf string) {
			doWork(j, sf)
			wg.Done()
		}(job, testFile)
	}

	wg.Wait()
}

func doWork(job *progress.Job, testFile string) {

	var wg *waitgroup.WaitGroup = waitgroup.New()

	numChunks := math.Ceil(magnitudes[testFile] / chunkSize)
	for i := 0; i <= int(numChunks); i++ {
		wg.Add(1)
		go func() {
			ct := ctr.Add()
			task, _ := job.GetTask(chunkName(ct))
			task.Start()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(7000)))
			task.End()
			job.UpdateBar()
			wg.Done()
		}()
	}

	wg.Wait()
}

func chunkName(i int) string {
	return fmt.Sprintf("chunk_%v", i)
}
