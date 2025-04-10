package rabbitmq

import "github.com/natanro/simlador-credito/motor-simulacao/infra"

type (
	simulationQueue struct {
		queue infra.Queue
	}

	SimulationQueue interface {
		Publish(msg *infra.QueueMessage) error
	}
)

func NewSimulationQueue(queue infra.Queue) SimulationQueue {
	return &simulationQueue{
		queue: queue,
	}
}

func (s *simulationQueue) Publish(msg *infra.QueueMessage) error {
	return s.queue.Publish(msg)
}
