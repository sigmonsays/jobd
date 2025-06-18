package schedule

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func NewScheduler() *Scheduler {
	c := cron.New()
	return &Scheduler{
		c:     c,
		eids:  make(map[int]*Job, 0),
		names: make(map[string]*Job, 0),
	}
}

type Scheduler struct {
	mx sync.Mutex

	c    *cron.Cron
	Jobs []*Job

	// map of jobs by eid (after scheduling occurs)
	eids map[int]*Job

	// map of jobs by name
	names map[string]*Job
}

func (me *Scheduler) Start() {
	me.c.Start()
}

func (me *Scheduler) Stop() context.Context {
	return me.c.Stop()
}

func (me *Scheduler) FindJobByName(jid string) (*Job, error) {
	me.mx.Lock()
	defer me.mx.Unlock()

	job, found := me.names[jid]
	if !found {
		return nil, fmt.Errorf("job not found: %s", jid)
	}

	return job, nil
}

type Next struct {
	Next, Prev string
}

func (me *Scheduler) FindNext(jid string) (*Next, error) {
	_, err := me.FindJobByName(jid)
	if err != nil {
		return nil, err
	}
	ret := &Next{}

	// find the jobs schedule
	entries := me.c.Entries()
	for _, entry := range entries {
		cj, ok := entry.Job.(*CronJob)
		if !ok {
			continue
		}
		if cj.Job.Id != jid {
			continue
		}

		ret.Prev = entry.Prev.Format(time.RFC3339)
		ret.Next = entry.Next.Format(time.RFC3339)

		break
	}

	return ret, nil
}

// small wrapper to provide an interface for the cron lib
// I can type assert later to find the exact job using the Entries method
type CronJob struct {
	Job *Job
}

func (me *CronJob) Run() {
	me.Job.Fn()
}

func (me *Scheduler) AddJob(j *Job) error {
	me.Jobs = append(me.Jobs, j)
	cj := &CronJob{j}
	eid, err := me.c.AddJob(j.Schedule, cj)
	slog.Info("Scheduled job", "jid", j.Id, "eid", eid)

	me.saveJob(j, int(eid))

	if err != nil {
		return err
	}
	return nil
}

func (me *Scheduler) saveJob(j *Job, eid int) error {
	me.mx.Lock()
	defer me.mx.Unlock()

	j.SetEid(eid)

	// save by eid
	me.eids[eid] = j

	// save by name
	me.names[j.Id] = j
	return nil
}

func NewJob() *Job {
	j := &Job{}
	return j
}

type Job struct {
	mx       sync.Mutex
	Eid      int
	Id       string
	Schedule string
	Fn       func() error `json:"-"`
	Running  bool
}

func (me *Job) SetEid(eid int) {
	me.Eid = eid
}

func (me *Job) SetRunning(running bool) {
	me.mx.Lock()
	defer me.mx.Unlock()
	me.Running = running
}

func (me *Job) Run() error {
	me.SetRunning(true)
	defer me.SetRunning(false)
	return me.Fn()
}
