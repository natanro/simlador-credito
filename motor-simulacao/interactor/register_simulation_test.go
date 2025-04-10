package interactor_test

import (
	"errors"
	"testing"

	"github.com/natanro/simlador-credito/motor-simulacao/entity"
	"github.com/natanro/simlador-credito/motor-simulacao/infra"
	"github.com/natanro/simlador-credito/motor-simulacao/interactor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para a fila de simulação
type MockSimulationQueue struct {
	mock.Mock
}

func (m *MockSimulationQueue) Publish(message *infra.QueueMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockQueue := new(MockSimulationQueue)

	register := interactor.NewSimulationRegister(mockRepo, mockQueue)

	simulation := &entity.Simulation{
		RequestedAmount: 10000.00,
		Installments:    12,
		Age:            30,
	}

	mockRepo.On("Create", mock.Anything).Return(1, nil)
	mockQueue.On("Publish", mock.Anything).Return(nil)

	err := register.Create(simulation)

	assert.NoError(t, err)
	assert.Equal(t, 1, simulation.ID)
	assert.Equal(t, entity.SimulationStatusCreated, simulation.Status)
	mockRepo.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

func TestCreate_NilSimulation(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockQueue := new(MockSimulationQueue)

	register := interactor.NewSimulationRegister(mockRepo, mockQueue)

	err := register.Create(nil)

	assert.Error(t, err)
	assert.Equal(t, "simulation cannot be nil", err.Error())
}

func TestCreate_CreateError(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockQueue := new(MockSimulationQueue)

	register := interactor.NewSimulationRegister(mockRepo, mockQueue)

	simulation := &entity.Simulation{
		RequestedAmount: 10000.00,
		Installments:    12,
		Age:            30,
	}

	mockRepo.On("Create", mock.Anything).Return(0, errors.New("create error"))

	err := register.Create(simulation)

	assert.Error(t, err)
	assert.Equal(t, "error creating simulation: create error", err.Error())
	mockRepo.AssertExpectations(t)
	mockQueue.AssertNotCalled(t, "Publish", mock.Anything)
}

func TestCreate_PublishError(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockQueue := new(MockSimulationQueue)

	register := interactor.NewSimulationRegister(mockRepo, mockQueue)

	simulation := &entity.Simulation{
		RequestedAmount: 10000.00,
		Installments:    12,
		Age:            30,
	}

	mockRepo.On("Create", mock.Anything).Return(1, nil)
	mockQueue.On("Publish", mock.Anything).Return(errors.New("publish error"))

	err := register.Create(simulation)

	assert.Error(t, err)
	assert.Equal(t, "error publishing simulation: publish error", err.Error())
	assert.Equal(t, 1, simulation.ID)
	assert.Equal(t, entity.SimulationStatusCreated, simulation.Status)
	mockRepo.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}
