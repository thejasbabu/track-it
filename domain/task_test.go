package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thejasbabu/track-it/domain"
)

func TestTaskCompletenessForDailyTask(t *testing.T) {
	currentTime := time.Date(2016, time.August, 15, 0, 0, 0, 0, time.Local)
	initTask := domain.NewTask(currentTime)
	assert.False(t, initTask.IsComplete(currentTime))

	task := domain.NewTask(currentTime)
	task.SetInterval("D")

	assert.True(t, task.Done(currentTime).IsComplete(currentTime))
	assert.False(t, task.UnDone(currentTime.Add(5*time.Minute)).IsComplete(currentTime))
	assert.False(t, task.Done(currentTime.AddDate(0, 0, 1)).IsComplete(currentTime))
	assert.False(t, task.Done(currentTime.AddDate(0, 0, -1)).IsComplete(currentTime))
}

func TestTaskCompletenessForWeeklyTask(t *testing.T) {
	currentTime := time.Date(2016, time.August, 15, 0, 0, 0, 0, time.Local)
	task := domain.NewTask(currentTime)
	task.SetInterval("W")

	assert.True(t, task.Done(currentTime).IsComplete(currentTime))
	assert.True(t, task.Done(currentTime.AddDate(0, 0, 1)).IsComplete(currentTime))
	assert.False(t, task.Done(currentTime.AddDate(0, 0, 8)).IsComplete(currentTime))
	assert.False(t, task.Done(currentTime.AddDate(0, 0, -8)).IsComplete(currentTime))
}

func TestTaskCompletenessForMonthlyTask(t *testing.T) {
	currentTime := time.Date(2016, time.August, 15, 0, 0, 0, 0, time.Local)
	task := domain.NewTask(currentTime)
	task.SetInterval("M")
	assert.False(t, task.Done(currentTime.AddDate(0, 1, 0)).IsComplete(currentTime))
	assert.True(t, task.Done(currentTime.AddDate(0, 0, -1)).IsComplete(currentTime))
	assert.True(t, task.Done(currentTime.AddDate(0, 0, -8)).IsComplete(currentTime))
	assert.False(t, task.Done(currentTime.AddDate(0, -1, 0)).IsComplete(currentTime))
}
