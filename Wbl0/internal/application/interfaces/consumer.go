package interfaces

// MessageConsumer определяет интерфейс для потребления сообщений из Kafka
type MessageConsumer interface {
	Start() error

	Stop() error

	Close() error
}
