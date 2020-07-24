package task_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
	util "github.com/thejasbabu/track-it/util/testlib/mocks"
)

func TestRepositoryWriteToBeSuccessful(t *testing.T) {
	badgerOps := new(util.Badger)
	repo := task.NewBadgerRepository(badgerOps)
	task := domain.Task{ID: "test1234", Description: "testing"}

	badgerOps.On("Update", []byte("test1234"), mock.Anything).Return(nil)
	err := repo.Write(task)

	require.NoError(t, err, "expected no error")
	mock.AssertExpectationsForObjects(t, badgerOps)
}

func TestRepositoryWriteToBeFailure(t *testing.T) {
	badgerOps := new(util.Badger)
	repo := task.NewBadgerRepository(badgerOps)
	task := domain.Task{ID: "test1234", Description: "testing"}

	badgerOps.On("Update", []byte("test1234"), mock.Anything).Return(errors.New("error"))
	err := repo.Write(task)

	assert.EqualError(t, err, "error")
	mock.AssertExpectationsForObjects(t, badgerOps)
}
