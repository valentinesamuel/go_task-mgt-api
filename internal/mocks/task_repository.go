package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) Get(id uint) (*models.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) List() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task *models.Task) (*models.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(id uint) (*models.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		// Return nil task and error for failure cases
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}
