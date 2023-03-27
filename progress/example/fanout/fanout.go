//nolint
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rockwell-uk/csync/waitgroup"
	"github.com/rockwell-uk/go-progress/progress"
	"github.com/rockwell-uk/uiprogress"
)

func main() {

	var jobName string = "testFanout"

	uiprogress.Start()

	var numFiles int = 10

	var tasks []*progress.Task
	for i := 0; i < numFiles; i++ {
		tasks = append(tasks, &progress.Task{
			ID:        fmt.Sprintf("file_%v", i),
			Magnitude: 1,
		})
	}

	job := progress.NewJob(jobName, numFiles)
	job.AddTasks(tasks)
	defer job.End(true)

	job.Start()

	fanout(job)

	uiprogress.Stop()

	fmt.Printf("%v ended\n", jobName)
}

func fanout(job *progress.Job) {

	var wg *waitgroup.WaitGroup = waitgroup.New()

	for taskName, task := range job.Tasks {

		wg.Add(1)
		go func(t *progress.Task, tn string) {
			t.Start()
			doWork(tn)
			t.End()
			wg.Done()
		}(task, taskName)
	}

	wg.Wait()
}

func doWork(l string) {

	var wg *waitgroup.WaitGroup = waitgroup.New()

	var numChunks int = 10
	var chunkNames []string

	for i := 0; i < numChunks; i++ {
		chunkNames = append(chunkNames, chunkName(i))
	}

	bar := uiprogress.AddBar(len(chunkNames)).AppendElapsed().PrependCompleted()
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return l
	})

	for i := range chunkNames {

		wg.Add(1)

		go func(b *uiprogress.Bar, k int) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
			b.Incr()
			wg.Done()
		}(bar, i)
	}

	wg.Wait()
}

func chunkName(i int) string {
	return fmt.Sprintf("chunk_%v", i)
}
