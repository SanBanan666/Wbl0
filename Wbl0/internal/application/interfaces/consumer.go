package interfaces

// MessageConsumer определяет интерфейс для потребления сообщений из Kafka
type MessageConsumer interface {
	// Start начинает потребление сообщений
	Start() error

	// Stop останавливает потребление сообщений
	Stop() error

	// Close закрывает соединение с Kafka
	Close() error
}
