package amqp

import (
	"fmt"
	"log"

	"github.com/natanro/simlador-credito/motor-simulacao/entity"
	"github.com/natanro/simlador-credito/motor-simulacao/infra"
	"github.com/natanro/simlador-credito/motor-simulacao/interactor"
)

type (
	simulationAmqpHandler struct {
		simulationProcessor interactor.SimulationProcessor
	}

	SimulationAmqpHandler interface {
		Notify(*infra.QueueMessage) error
	}
)

func NewSimulationAmqpHandler(simulationProcessor interactor.SimulationProcessor) SimulationAmqpHandler {
	s := &simulationAmqpHandler{
		simulationProcessor: simulationProcessor,
	}

	return s
}

func (s *simulationAmqpHandler) Notify(msg *infra.QueueMessage) error {
	if msg == nil {
		return fmt.Errorf("invalid message")
	}

	log.Println("Received message: ", msg.ID)
	
	simulation, ok := msg.Message.(*entity.Simulation)
	if !ok {
		return fmt.Errorf("invalid message type")
	}
	if simulation == nil {
		return nil
	}
	if simulation.Status != entity.SimulationStatusCreated {
		return nil
	}
	return s.simulationProcessor.Consume(simulation)
}
