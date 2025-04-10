package adapter

import (
	"github.com/natanro/simlador-credito/motor-simulacao/entity"
	"github.com/natanro/simlador-credito/motor-simulacao/infra"
	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db/model"
)

func SimulationModelToEntity(s *model.Simulation) *entity.Simulation {
	return &entity.Simulation{
		ID:              s.ID,
		RequestedAmount: s.RequestedAmount,
		Installments:    s.Installments,
		Status:          entity.SimulationStatus(s.Status),
		Age:             s.Age,
		AnnualRate:      s.AnnualRate,
		MonthlyRate:     s.MonthlyRate,
		MonthlyPayment:  s.MonthlyPayment,
		TotalAmount:     s.TotalAmount,
	}
}

func SimulationEntityToModel(s *entity.Simulation) *model.Simulation {
	return &model.Simulation{
		ID:              s.ID,
		RequestedAmount: s.RequestedAmount,
		Installments:    s.Installments,
		Status:          string(s.Status),
		Age:             s.Age,
		AnnualRate:      s.AnnualRate,
		MonthlyRate:     s.MonthlyRate,
		MonthlyPayment:  s.MonthlyPayment,
		TotalAmount:     s.TotalAmount,
	}
}

func SimulationEntityToQueueMessage(s *entity.Simulation) *infra.QueueMessage {
	return &infra.QueueMessage{
		ID:      s.ID,
		Message: s,
	}
}
