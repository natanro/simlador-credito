package interactor

import (
	"errors"
	"fmt"
	"log"

	"github.com/natanro/simlador-credito/motor-simulacao/adapter"
	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db"
	"github.com/natanro/simlador-credito/motor-simulacao/datasource/rabbitmq"
	"github.com/natanro/simlador-credito/motor-simulacao/entity"
)

type (
	simulationRegister struct {
		simulationRepository db.SimulationRepository
		simulationQueue      rabbitmq.SimulationQueue
	}

	SimulationRegister interface {
		Create(simulation *entity.Simulation) error
	}
)

func NewSimulationRegister(simulationRepository db.SimulationRepository, simulationQueue rabbitmq.SimulationQueue) SimulationRegister {
	return &simulationRegister{
		simulationRepository: simulationRepository,
		simulationQueue:      simulationQueue,
	}
}

func (s *simulationRegister) Create(simulation *entity.Simulation) error {
	if simulation == nil {
		return errors.New("simulation cannot be nil")
	}

	simulation.Status = entity.SimulationStatusCreated
	simulationID, err := s.simulationRepository.Create(adapter.SimulationEntityToModel(simulation))
	if err != nil {
		return fmt.Errorf("error creating simulation: %w", err)
	}

	simulation.ID = simulationID
	if err := s.simulationQueue.Publish(adapter.SimulationEntityToQueueMessage(simulation)); err != nil {
		return fmt.Errorf("error publishing simulation: %w", err)
	}

	log.Println("simulation created successfully")

	return nil
}
