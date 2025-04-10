package interactor_test

import (
	"errors"
	"testing"

	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db/model"
	"github.com/natanro/simlador-credito/motor-simulacao/entity"
	"github.com/natanro/simlador-credito/motor-simulacao/interactor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para o repositório de simulação
type MockSimulationRepository struct {
	mock.Mock
}

func (m *MockSimulationRepository) Create(simulation *model.Simulation) (int, error) {
	args := m.Called(simulation)
	return args.Int(0), args.Error(1)
}

func (m *MockSimulationRepository) Update(simulation *model.Simulation) error {
	args := m.Called(simulation)
	return args.Error(0)
}

// Mock para a estratégia de taxa
type MockRateStrategy struct {
	mock.Mock
}

func (m *MockRateStrategy) GetRateByAge(age int) (float64, error) {
	args := m.Called(age)
	return args.Get(0).(float64), args.Error(1)
}

func TestConsume_Success(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockRateStrategy := new(MockRateStrategy)

	processor := interactor.NewSimulationProcessor(mockRepo, mockRateStrategy)

	simulation := &entity.Simulation{
		Age:             30,
		RequestedAmount: 10000,
		Installments:    12,
	}

	mockRateStrategy.On("GetRateByAge", simulation.Age).Return(0.05, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	err := processor.Consume(simulation)

	assert.NoError(t, err)
	assert.Equal(t, 0.05, simulation.AnnualRate)
	assert.Equal(t, 0.004166666666666667, simulation.MonthlyRate) // 0.05 / 12 -- precision
	assert.NotZero(t, simulation.TotalAmount)
	assert.NotZero(t, simulation.MonthlyPayment)
	mockRepo.AssertExpectations(t)
	mockRateStrategy.AssertExpectations(t)
}

func TestConsume_GetRateByAgeError(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockRateStrategy := new(MockRateStrategy)

	processor := interactor.NewSimulationProcessor(mockRepo, mockRateStrategy)

	simulation := &entity.Simulation{
		Age:             30,
		RequestedAmount: 10000,
		Installments:    12,
	}

	mockRateStrategy.On("GetRateByAge", simulation.Age).Return(0.0, errors.New("rate error"))

	err := processor.Consume(simulation)

	assert.Error(t, err)
	assert.Equal(t, "rate error", err.Error())
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestConsume_UpdateError(t *testing.T) {
	mockRepo := new(MockSimulationRepository)
	mockRateStrategy := new(MockRateStrategy)

	processor := interactor.NewSimulationProcessor(mockRepo, mockRateStrategy)

	simulation := &entity.Simulation{
		Age:             30,
		RequestedAmount: 10000,
		Installments:    12,
	}

	mockRateStrategy.On("GetRateByAge", simulation.Age).Return(0.05, nil)
	mockRepo.On("Update", mock.Anything).Return(errors.New("update error"))

	err := processor.Consume(simulation)

	assert.Error(t, err)
	assert.Equal(t, "error updating simulation: update error", err.Error())
}
