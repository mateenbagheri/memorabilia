package schedule

import (
	"time"
)

type CronJobExpr string

type CronjobRepository interface {
	ScheduleIntervalJob(timeExpr string, job func()) (jobID string, err error)
	RemoveJob(jobID string) error
	Start() error
	Stop() error
	ListJobs() ([]JobDetails, error)
}

type JobDetails struct {
	JobId   string
	NextRun time.Time
	LastRun time.Time
}
