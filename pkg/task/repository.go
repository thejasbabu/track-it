package task

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/util"
)

// Repository is an interface for Task storage
type Repository interface {
	Write(task domain.Task) error
	FetchAll() ([]domain.Task, error)
	Delete(task domain.Task) error
}

// BadgerRepository is a wrapper over task storage using badger
type BadgerRepository struct {
	badgerOp util.Badger
}

// NewBadgerRepository returns a new BadgerRepository
func NewBadgerRepository(badgerOp util.Badger) BadgerRepository {
	return BadgerRepository{badgerOp: badgerOp}
}

// Writes the task to the badgerDB
func (r BadgerRepository) Write(t domain.Task) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(t)
	if err != nil {
		return fmt.Errorf("badgerrepository: encode task: %w", err)
	}
	return r.badgerOp.Update([]byte(t.ID), buffer.Bytes())
}

// FetchAll returns all tasks
func (r BadgerRepository) FetchAll() ([]domain.Task, error) {
	tasks := []domain.Task{}
	err := r.badgerOp.Read(func(val []byte) error {
		var task domain.Task
		buffer := bytes.NewBuffer(val)
		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&task); err != nil {
			return err
		}
		tasks = append(tasks, task)
		return nil
	})
	return tasks, err
}

func (r BadgerRepository) Delete(task domain.Task) error {
	return r.badgerOp.Delete([]byte(task.ID))
}
