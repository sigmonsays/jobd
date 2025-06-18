package app

import (
	"github.com/sigmonsays/jobd/config"
	"github.com/sigmonsays/jobd/schedule"
)

type Api struct {
	*Context
}
type Context struct {
	*config.AppConfig
	JobConfig *schedule.JobConfig
	Scheduler *schedule.Scheduler
	Executor  *schedule.Executor
}
