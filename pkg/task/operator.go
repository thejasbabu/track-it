package task

import (
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/util"
)

// Operator Handles the Task CRUD operations
type Operator struct {
	repository Repository
	identifier util.Identifier
	clock      util.Clock
}

// NewOperator returns a new Task Operator
func NewOperator(repo Repository, identifier util.Identifier, clock util.Clock) Operator {
	return Operator{repository: repo, identifier: identifier, clock: clock}
}

// Add a task
func (o Operator) Add(task domain.Task) error {
	task.ID = o.identifier.Generate()
	return o.repository.Write(task)
}

// MarkAsComplete marks the task as completed
func (o Operator) MarkAsDone(task domain.Task) error {
	return o.repository.Write(task.Done(o.clock.CurrentTime()))
}

func (o Operator) MarkAsUnDone(task domain.Task) error {
	return o.repository.Write(task.UnDone(o.clock.CurrentTime()))
}

func (o Operator) DeleteTask(task domain.Task) error {
	return o.repository.Delete(task)
}

// GetTasks returns all the tasks
// TODO: Paginate responses, don't fetch all
func (o Operator) GetTasks() ([]domain.Task, error) {
	return o.repository.FetchAll()
}
