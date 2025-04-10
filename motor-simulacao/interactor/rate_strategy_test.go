package interactor_test

import (
	"errors"
	"testing"

	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db"
	"github.com/natanro/simlador-credito/motor-simulacao/interactor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para o repositório de parâmetros de simulação
type MockSimulationParamRepository struct {
	mock.Mock
}

func (m *MockSimulationParamRepository) FindParams() ([]db.Document, error) {
	args := m.Called()
	return args.Get(0).([]db.Document), args.Error(1)
}

func TestGetRateByAge_Success(t *testing.T) {
	mockRepo := new(MockSimulationParamRepository)
	rateStrategy := interactor.NewRateStrategy(mockRepo)

	params := []db.Document{
		{Class: "25-", Rate: 0.05},
		{Class: "26-40", Rate: 0.03},
		{Class: "41-60", Rate: 0.02},
		{Class: "61+", Rate: 0.04},
	}

	mockRepo.On("FindParams").Return(params, nil)

	rate, err := rateStrategy.GetRateByAge(30)

	assert.NoError(t, err)
	assert.Equal(t, 0.03, rate) // Espera-se que a Rate para a idade 30 seja 0.03
	mockRepo.AssertExpectations(t)
}

func TestGetRateByAge_FindParamsError(t *testing.T) {
	mockRepo := new(MockSimulationParamRepository)
	rateStrategy := interactor.NewRateStrategy(mockRepo)

	mockRepo.On("FindParams").Return([]db.Document{}, errors.New("database error"))

	rate, err := rateStrategy.GetRateByAge(30)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Equal(t, float64(-1), rate) // Espera-se que a Rate seja -1
	mockRepo.AssertExpectations(t)
}
