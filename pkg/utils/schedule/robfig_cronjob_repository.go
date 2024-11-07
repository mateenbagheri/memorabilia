package schedule

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mateenbagheri/memorabilia/pkg/utils/validation"
	"github.com/robfig/cron/v3"
)

var _ CronjobRepository = (*RobfigCronjobRepository)(nil)

var (
	ErrScheduleFailed = errors.New("failed to schedule job")
	ErrJobNotFound    = errors.New("job id not found")
)

type RobfigCronjobRepository struct {
	scheduler *cron.Cron
	mu        sync.RWMutex
	jobs      map[string]cron.EntryID
}

func newRobfigCronjobRepository() *RobfigCronjobRepository {
	return &RobfigCronjobRepository{
		scheduler: cron.New(),
		jobs:      make(map[string]cron.EntryID),
	}
}

// Singleton instance of RobfigCronjobRepository and sync.Once to ensure it's only initialized once
var (
	instance *RobfigCronjobRepository
	once     sync.Once
)

// GetRobfigSchedulerInstance returns the singleton instance of RobfigCronjobRepository.
func GetRobfigSchedulerInstance() *RobfigCronjobRepository {
	once.Do(func() {
		instance = newRobfigCronjobRepository()
	})
	return instance
}

func (cj *RobfigCronjobRepository) ScheduleIntervalJob(timeExpr string, job func()) (jobId string, err error) {
	cj.mu.Lock()
	defer cj.mu.Unlock()

	err = validation.ValidateJobTimeFormat(timeExpr)
	if err != nil {
		return "", err
	}

	entryId, err := cj.scheduler.AddFunc(fmt.Sprintf("@every %s", timeExpr), job)
	if err != nil {
		return "", fmt.Errorf("%w, %v", ErrScheduleFailed, err)
	}

	jobId = fmt.Sprintf("%d", entryId) // entryId equals jobId
	cj.jobs[jobId] = entryId
	return jobId, nil
}

func (cj *RobfigCronjobRepository) RemoveJob(jobId string) error {
	cj.mu.Lock()
	defer cj.mu.Unlock()

	entryId, exists := cj.jobs[jobId]
	if !exists {
		return fmt.Errorf("%w, jobId: %s", ErrJobNotFound, jobId)
	}

	cj.scheduler.Remove(entryId)
	delete(cj.jobs, jobId)
	return nil
}

func (cj *RobfigCronjobRepository) Start() error {
	cj.scheduler.Start()
	return nil
}

func (cj *RobfigCronjobRepository) Stop() error {
	cj.scheduler.Stop()
	return nil
}

func (cj *RobfigCronjobRepository) ListJobs() ([]JobDetails, error) {
	cj.mu.RLock()
	defer cj.mu.RUnlock()

	var jobList []JobDetails
	for jobId, entryId := range cj.jobs {
		entry := cj.scheduler.Entry(entryId)
		jobList = append(jobList, JobDetails{
			JobId:   jobId,
			NextRun: entry.Next,
			LastRun: entry.Prev,
		})
	}
	return jobList, nil
}
