package domain

import (
	"strings"
	"time"
)

// Interval represents the time grain for task completion
type Interval string

const (
	// NONE interval for one-time task
	NONE Interval = "NONE"
	// DAILY interval for Daily tasks
	DAILY Interval = "DAILY"
	// WEEKLY interval for Weekly tasks
	WEEKLY Interval = "WEEKLY"
	// MONTHLY interval for Monthly tasks
	MONTHLY Interval = "MONTHLY"
)

// Task represents a Task to be completed
type Task struct {

	// Uniquely representing the Task
	ID string

	// Description of the task
	Description string

	// Tags are used to identify related tasks
	Tags []string

	// RepeatInterval for checking task completion
	RepeatInterval Interval

	// Streak represents the number of days continously completed
	Streak int32

	// streak represents the number of days continously completed
	Records []Record
}

// INITSTREAK set to zero
const INITSTREAK int32 = 0

// NewTask returns a new Task obj
func NewTask(t time.Time) Task {
	record := Record{Status: UNDONE, LastUpdatedTime: t}
	return Task{Streak: INITSTREAK, Records: []Record{record}}
}

// Assign an ID to the task
func (t Task) Assign(id string) Task {
	t.ID = id
	return t
}

func (t *Task) SetInterval(interval string) {
	switch {
	case strings.HasPrefix(interval, "D"):
		t.RepeatInterval = DAILY
	case strings.HasPrefix(interval, "W"):
		t.RepeatInterval = WEEKLY
	case strings.HasPrefix(interval, "M"):
		t.RepeatInterval = MONTHLY
	default:
		t.RepeatInterval = NONE
	}
}

func (t Task) Done(term time.Time) Task {
	if !t.IsComplete(term) {
		t.Records = append(t.Records, Record{Status: DONE, LastUpdatedTime: term})
		t.Streak++
	}
	return t
}

func (t Task) UnDone(term time.Time) Task {
	if t.IsComplete(term) {
		t.Records = append(t.Records, Record{Status: UNDONE, LastUpdatedTime: term})
		t.Streak--
	}
	return t
}

// IsComplete checks if the task is completed on not
func (t Task) IsComplete(term time.Time) bool {
	switch t.RepeatInterval {
	case NONE:
		return len(t.Records) > 0 && t.Records[len(t.Records)-1].Status == DONE
	default:
		return len(t.Records) > 0 && checkCompleteness(t.Records, t.RepeatInterval, term)
	}
}

func checkCompleteness(recs []Record, intervalType Interval, term time.Time) bool {
	switch intervalType {
	case DAILY:
		lastDayRec := recs[len(recs)-1]
		return lastDayRec.Status == DONE && term.YearDay() == lastDayRec.LastUpdatedTime.YearDay() &&
			term.Year() == lastDayRec.LastUpdatedTime.Year()
	case WEEKLY:
		lastDayRec := recs[len(recs)-1]
		termYear, termWeek := term.ISOWeek()
		lastUpdatedYear, lastUpdatedWeek := lastDayRec.LastUpdatedTime.ISOWeek()
		return lastDayRec.Status == DONE && termWeek == lastUpdatedWeek && termYear == lastUpdatedYear
	case MONTHLY:
		lastDayRec := recs[len(recs)-1]
		return lastDayRec.Status == DONE && term.Month() == lastDayRec.LastUpdatedTime.Month() &&
			term.Year() == lastDayRec.LastUpdatedTime.Year()
	default:
		return false
	}
}

// Status represents the status of the task
type Status string

const (
	// DONE status means that the task was completed within the Interval
	DONE Status = "DONE"
	// UNDONE status means that the task was not completed within the Interval
	UNDONE Status = "UNDONE"
)

// Record tracks the tasks along with it's status
type Record struct {

	// Status of the task
	Status Status

	// LastUpdatedTime represents the task status update time
	LastUpdatedTime time.Time
}
