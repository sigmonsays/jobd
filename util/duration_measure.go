package util

import "time"

// simple helper function to measure start/stop times
func NewDurationMeasure(start bool) *DurationMeasure {
	ret := &DurationMeasure{}
	if start {
		ret.Start = time.Now()
	}
	return ret
}

type DurationMeasure struct {
	Start time.Time
	Stop  time.Time
}

func (me *DurationMeasure) StartTimer() *DurationMeasure {
	me.Start = time.Now()
	return me
}

func (me *DurationMeasure) StopTimer() *DurationMeasure {
	me.Stop = time.Now()
	return me
}

func (me *DurationMeasure) GetDuration() time.Duration {
	return me.Stop.Sub(me.Start)
}

func (me *DurationMeasure) GetDurationSec() int64 {
	d := me.GetDuration()
	return int64(d.Seconds())
}
