package infra

type (
	queue struct {
		name string
		queue chan *QueueMessage
		observers []Observer
		retries int
		retriesCount map[int]int
	}

	Observer interface {
		Notify(message *QueueMessage) error
	}

	Queue interface {
		Publish(message *QueueMessage) error
		RegisterObserver(observer Observer)
	}

	QueueMessage struct {
		ID int
		Message interface{}
	}
)

func NewQueue(queueName string, retries int) Queue {
	q := &queue{
		name: queueName,
		queue: make(chan *QueueMessage, 10000),
		retries: retries,
		retriesCount: make(map[int]int),
	}
	go q.runListener()
	return q
}

func (q *queue) Publish(message *QueueMessage) error {
	if _, ok := q.retriesCount[message.ID]; !ok {
		q.retriesCount[message.ID] = 0
	}
	if q.retriesCount[message.ID] < q.retries {
		q.retriesCount[message.ID] += 1
		q.queue <- message
	}
	return nil
}

func (q *queue) RegisterObserver(observer Observer) {
	q.observers = append(q.observers, observer)
}

func (q *queue) runListener() {
	for message := range q.queue {
		for _, observer := range q.observers {
			if err := observer.Notify(message); err != nil {
				q.Publish(message)
			}
		}
	}
}

func (q *queue) countRetries(id int) int {
	return q.retriesCount[id]
}

func (q *queue) getObservers() []Observer {
	return q.observers
}