package infra

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para o observador
type MockObserver struct {
	mock.Mock
}

func (m *MockObserver) Notify(message *QueueMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestPublish_Success(t *testing.T) {
	queue := &queue{
		name:         "test_queue",
		queue:        make(chan *QueueMessage, 10000),
		observers:    []Observer{},
		retries:      10,
		retriesCount: make(map[int]int),
	}
	message := &QueueMessage{ID: 1, Message: "Test message"}

	err := queue.Publish(message)

	assert.NoError(t, err)
	assert.Equal(t, 1, queue.countRetries(message.ID)) // Deve ter tentado apenas uma vez
}

func TestPublish_MaxRetries(t *testing.T) {
	mockObserver := &MockObserver{}
	queue := &queue{
		name:         "test_queue",
		queue:        make(chan *QueueMessage, 10000),
		observers:    []Observer{mockObserver},
		retries:      2, // Define o limite de retries para 2
		retriesCount: make(map[int]int),
	}
	message := &QueueMessage{ID: 1, Message: "Test message"}

	mockObserver.On("Notify", message).Return(errors.New("error")).Twice()
	go queue.runListener()

	// Publicar duas vezes para exceder o limite de retries
	queue.Publish(message)

	time.Sleep(100 * time.Millisecond) // Espera um pouco para garantir que a mensagem seja processada

	assert.Equal(t, 2, queue.countRetries(message.ID)) // Deve ter tentado duas vezes
}

func TestRunListener_NotifySuccess(t *testing.T) {
	queue := &queue{
		name:         "test_queue",
		queue:        make(chan *QueueMessage),
		observers:    []Observer{},
		retries:      10,
		retriesCount: make(map[int]int),
	}
	mockObserver := new(MockObserver)
	queue.RegisterObserver(mockObserver)

	message := &QueueMessage{ID: 1, Message: "Test message"}

	mockObserver.On("Notify", message).Return(nil).Once()

	go queue.runListener()

	queue.Publish(message) // Publicar a mensagem

	time.Sleep(100 * time.Millisecond)

	// Verificar se o observador foi chamado
	mockObserver.AssertExpectations(t)
}
