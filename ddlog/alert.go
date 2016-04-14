package ddlog

import (
	"fmt"
	"time"
)

type Alert struct {
	Start time.Time
	End   time.Time
	Count int
}

func (a *Alert) String() string {
	if a.IsDone() {
		return fmt.Sprintf("[Alert Done] at %s duration=%v\n", a.End, a.End.Sub(a.Start))
	}
	return fmt.Sprintf("High traffic generated an alert - hits = %v, triggered at %s", a.Count, a.Start.Format(DefaultTimeFormat))
}

func (a *Alert) Copy() *Alert {
	return &Alert{Start: a.Start, End: a.End, Count: a.Count}
}

func (a *Alert) Complete(at time.Time) {
	a.End = at
}

func (a *Alert) IsDone() bool {
	return !a.End.IsZero()
}
