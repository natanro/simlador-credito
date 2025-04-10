package entity

import "github.com/natanro/simlador-credito/motor-simulacao/datasource/db/model"

type (
	SimulationStatus string
	Simulation       struct {
		ID              int              `json:"id"`
		RequestedAmount float64          `json:"requested_amount"`
		Installments    int              `json:"installments"`
		Status          SimulationStatus `json:"status"`
		Age             int              `json:"age"`
		AnnualRate      float64          `json:"annual_rate"`
		MonthlyRate     float64          `json:"monthly_rate"`
		MonthlyPayment  float64          `json:"monthly_payment"`
		TotalAmount    float64          `json:"total_amount"`
	}
)

const (
	SimulationStatusCreated   SimulationStatus = "CREATED"
	SimulationStatusProcessed SimulationStatus = "PROCESSED"
)

func (s *Simulation) ToModel() model.Simulation {
	return model.Simulation{
		ID:              s.ID,
		RequestedAmount: s.RequestedAmount,
		Installments:    s.Installments,
		Status:          string(s.Status),
		Age:             s.Age,
		AnnualRate:      s.AnnualRate,
		MonthlyRate:     s.MonthlyRate,
		MonthlyPayment:  s.MonthlyPayment,
	}
}
