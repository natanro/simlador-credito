package interactor

import (
	"fmt"
	"log"
	"math"

	"github.com/natanro/simlador-credito/motor-simulacao/adapter"
	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db"
	"github.com/natanro/simlador-credito/motor-simulacao/entity"
)

type (
	simulationProcessor struct {
		simulationRepository db.SimulationRepository
		rateStrategy RateStrategy
	}

	SimulationProcessor interface {
		Consume(simulation *entity.Simulation) error
	}
)

func NewSimulationProcessor(simulationRepository db.SimulationRepository, rateStrategy RateStrategy) SimulationProcessor {
	return &simulationProcessor{
		simulationRepository: simulationRepository,
		rateStrategy: rateStrategy,
	}
}

func (s *simulationProcessor) Consume(simulation *entity.Simulation) error {
	log.Println("Consuming simulation: ", simulation)
	
	annualRate, err := s.rateStrategy.GetRateByAge(simulation.Age)
	if err != nil {
		log.Println("Erro ao obter taxa: ", err)
		return err
	}

	simulation.AnnualRate = annualRate
	simulation.MonthlyRate = annualRate / 12

	log.Println("Annual rate: ", annualRate, "Monthly rate: ", simulation.MonthlyRate)

	numerator := simulation.RequestedAmount * (1 + simulation.MonthlyRate)
	denominator := (1 / math.Pow(1 + simulation.MonthlyRate, float64(simulation.Installments)))
	simulation.TotalAmount = numerator / denominator
	simulation.MonthlyPayment = simulation.TotalAmount / float64(simulation.Installments) // precision?

	log.Println("Monthly payment calculated: ", simulation.MonthlyPayment)

	simulation.Status = entity.SimulationStatusProcessed
	if err := s.simulationRepository.Update(adapter.SimulationEntityToModel(simulation)); err != nil {
		return fmt.Errorf("error updating simulation: %w", err)
	}

	log.Println("Simulation processed successfully")

	return nil
}