package task_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thejasbabu/track-it/domain"
	"github.com/thejasbabu/track-it/pkg/task"
	"github.com/thejasbabu/track-it/pkg/testlib/mocks"
	util "github.com/thejasbabu/track-it/util/testlib/mocks"
)

func TestOperatorAddFunctionToBeSuccess(t *testing.T) {
	repo := new(mocks.Repository)
	identifier := new(util.Identifier)
	clock := new(util.Clock)
	op := task.NewOperator(repo, identifier, clock)
	task := domain.NewTask(time.Now())
	task.ID = "test1234"
	identifier.On("Generate").Return("test1234")
	repo.On("Write", task.Assign("test1234")).Return(nil)

	err := op.Add(task)
	require.NoError(t, err, "expected no error")
	mock.AssertExpectationsForObjects(t, identifier, repo)
}

func TestOperatorAddFunctionToFail(t *testing.T) {
	repo := new(mocks.Repository)
	identifier := new(util.Identifier)
	clock := new(util.Clock)
	op := task.NewOperator(repo, identifier, clock)
	task := domain.NewTask(time.Now())
	task.ID = "test1234"

	identifier.On("Generate").Return("test1234")
	repo.On("Write", task.Assign("test1234")).Return(errors.New("error"))

	err := op.Add(task)
	assert.EqualError(t, err, "error")
	mock.AssertExpectationsForObjects(t, identifier, repo)
}
