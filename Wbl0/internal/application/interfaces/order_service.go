package interfaces

import "WbServis/Wbl0/internal/domain/entities"

// OrderService определяет интерфейс для бизнес-логики работы с заказами
type OrderService interface {
	// ProcessOrder обрабатывает новый заказ (сохраняет в БД и кэш)
	ProcessOrder(order *entities.Order) error

	// GetOrderByID получает заказ по ID (сначала из кэша, затем из БД)
	GetOrderByID(orderUID string) (*entities.Order, error)

	// RestoreCache восстанавливает кэш из базы данных при запуске
	RestoreCache() error

	// ProcessMessage обрабатывает сообщение из Kafka
	ProcessMessage(message []byte) error

	// Close закрывает сервис
	Close() error
}
