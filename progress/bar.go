//nolint:errcheck
package progress

func (j *Job) SetBar(v int) {
	if ShouldShowBar() {
		j.Bar.Set(v)
	}
}

func (j *Job) IncrBar() {
	if ShouldShowBar() {
		if j.Bar.CompletedPercent() < 100 {
			j.Bar.Incr()
		}
	}
}

func (j *Job) UpdateBar() {
	if ShouldShowBar() {
		done, _ := j.GetDone()
		j.Bar.Set(done)
	}
}

func (j *Job) EndBar() {
	if ShouldShowBar() {
		j.SetBar(int(j.Magnitude))
	}
}
