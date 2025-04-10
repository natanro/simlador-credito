package interactor

import (
	"errors"

	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db"
)

type RateStrategy interface {
	GetRateByAge(age int) (float64, error)
}

type rateStrategy struct {
	paramRepository db.SimulationParamRepository
}

func NewRateStrategy(paramRepository db.SimulationParamRepository) RateStrategy {
	return &rateStrategy{
		paramRepository: paramRepository,
	}
}

// GetRateByAge retrieves the interest rate based on the provided age by matching the age to predefined age classes.
// It returns the corresponding rate if a matching class is found, or an error if no suitable class exists.
func (t *rateStrategy) GetRateByAge(age int) (float64, error) {
	params, err := t.paramRepository.FindParams()
	if err != nil {
		return -1, err
	}

	for _, param := range params {
		if age < 25 {
			if param.Class == "25-" {
				return param.Rate, nil
			}
		} else if age >= 25 && age <= 40 {
			if param.Class == "26-40" {
				return param.Rate, nil
			}
		} else if age > 40 && age <= 60 {
			if param.Class == "41-60" {
				return param.Rate, nil
			}
		} else {
			if param.Class == "61+" {
				return param.Rate, nil
			}
		}
	}

	return 0, errors.New("classe não encontrada para a idade fornecida")
}
