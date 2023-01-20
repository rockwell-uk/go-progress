//nolint
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rockwell-uk/csync/waitgroup"
	"github.com/rockwell-uk/go-progress/progress"
	"github.com/rockwell-uk/uiprogress"
	"github.com/rockwell-uk/uiprogress/util/strutil"
)

var steps = []string{
	"downloading source",
	"installing deps",
	"compiling",
	"packaging",
	"seeding database",
	"deploying",
	"staring servers",
}

func main() {

	fmt.Println("apps: deployment started")
	uiprogress.Start()

	var wg *waitgroup.WaitGroup = waitgroup.New()

	wg.Add(1)
	go deploy("app1", wg)
	wg.Add(1)
	go deploy("app2", wg)
	wg.Add(1)
	go deploy("app3", wg)
	wg.Add(1)
	go deploy("app4", wg)
	wg.Add(1)
	go deploy("app5", wg)
	wg.Add(1)
	go deploy("app6", wg)

	wg.Wait()

	uiprogress.Stop()

	fmt.Println("apps: successfully deployed")
}

func deploy(app string, wg *waitgroup.WaitGroup) {

	var tasks []*progress.Task
	for _, step := range steps {
		tasks = append(tasks, &progress.Task{
			ID:        step,
			Magnitude: 1,
		})
	}

	job := progress.NewJob(app, len(steps))
	job.AddTasks(tasks)
	job.CalculateMagnitude()

	defer wg.Done()
	job.Bar = uiprogress.AddBar(len(steps)).AppendCompleted().PrependElapsed()
	job.Bar.Width = 50

	// Prepend the deploy step to the bar
	job.Bar.PrependFunc(func(b *uiprogress.Bar) string {
		return strutil.Resize(job.Name+": "+steps[b.Current()-1], 22)
	})
	job.Start()
	defer job.End(true)

	rand.Seed(500)
	for job.Bar.Incr() {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	}
}
