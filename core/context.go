package core

import (
	"github.com/sigmonsays/jobd/config"
	"github.com/sigmonsays/jobd/schedule"
)

type Context struct {
	*config.AppConfig
	JobConfig *schedule.JobConfig
	Scheduler *schedule.Scheduler
	Executor  *schedule.Executor
}
