package schedule

import (
	"fmt"
	"sync"

	"github.com/sigmonsays/jobd/job"
)

func NewJobConfig() *JobConfig {
	return &JobConfig{
		jobs: make(map[string]*job.JobSpec, 0),
	}
}

type JobConfig struct {
	mx   sync.Mutex
	jobs map[string]*job.JobSpec
}

func (me *JobConfig) ListJobs() []string {
	me.mx.Lock()
	defer me.mx.Unlock()
	ret := make([]string, 0)
	for jid, _ := range me.jobs {
		ret = append(ret, jid)
	}
	return ret
}
func (me *JobConfig) SetConfig(jid string, cfg *job.JobSpec) {
	me.mx.Lock()
	defer me.mx.Unlock()
	me.jobs[jid] = cfg
}

func (me *JobConfig) GetConfig(jid string) (*job.JobSpec, error) {
	me.mx.Lock()
	defer me.mx.Unlock()
	ret, found := me.jobs[jid]
	if !found {
		return nil, fmt.Errorf("job config not found: jid %s", jid)
	}
	return ret, nil
}
