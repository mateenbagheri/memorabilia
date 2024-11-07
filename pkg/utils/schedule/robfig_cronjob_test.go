package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRobfigCronjobRepository_StartStop(t *testing.T) {
	scheduler := GetRobfigSchedulerInstance()

	// Start the scheduler
	err := scheduler.Start()
	assert.NoError(t, err, "Expected no error when starting the scheduler")

	// Stop the scheduler
	err = scheduler.Stop()
	assert.NoError(t, err, "Expected no error when stopping the scheduler")
}

func TestRobfigCronjobRepository_RemoveJob(t *testing.T) {
	scheduler := GetRobfigSchedulerInstance()

	scheduler.Start()

	job := func() {}

	jobID, err := scheduler.ScheduleIntervalJob("1h", job)
	require.NoError(t, err, "Expected no error when scheduling a job")

	err = scheduler.RemoveJob(jobID)
	require.NoError(t, err, "Expected no error when removing the job")

	// Verify the job was removed
	err = scheduler.RemoveJob(jobID)
	assert.ErrorIs(t, err, ErrJobNotFound, "Expected job not found error when removing the job again")
}

func TestRobfigCronjobRepository_ScheduleIntervalJob(t *testing.T) {

}

func TestRobfigCronjobRepository_Listjobs(t *testing.T) {
	scheduler := GetRobfigSchedulerInstance()

	job1 := func() {}
	job2 := func() {}

	scheduler.Start()
	defer scheduler.Stop()

	job1Id, err := scheduler.ScheduleIntervalJob("30s", job1)
	require.NoError(t, err, "Expected no error while scheduling interval job number 1")

	job2Id, err := scheduler.ScheduleIntervalJob("20s", job2)
	require.NoError(t, err, "Expected no error while scheduling interval job number 2")

	jobs, err := scheduler.ListJobs()
	require.NoError(t, err, "Expected no error while getting job list")

	type Mapped struct {
		Exists  bool
		NextRun time.Time
		LastRun time.Time
	}
	jobIds := map[string]Mapped{
		job1Id: {Exists: false},
		job2Id: {Exists: false},
	}

	for _, job := range jobs {
		if _, exists := jobIds[job.JobId]; exists {
			jobIds[job.JobId] = Mapped{
				Exists:  true,
				NextRun: job.NextRun,
				LastRun: job.LastRun,
			}
		}
	}

	assert.True(t, jobIds[job1Id].Exists, "Expected a job with id of first created job to exist in ListJobs response")
	assert.True(t, jobIds[job2Id].Exists, "Expected a job with id of second created job to exist in ListJobs response")

	assert.Equal(t, jobIds[job1Id].LastRun, time.Time{}, "")
	assert.Equal(t, jobIds[job2Id].LastRun, time.Time{}, "")

	assert.NotEqual(t, jobIds[job1Id].NextRun, time.Time{})
	assert.NotEqual(t, jobIds[job2Id].NextRun, time.Time{})
}
