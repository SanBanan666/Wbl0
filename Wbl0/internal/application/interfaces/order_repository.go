package interfaces

import "WbServis/Wbl0/internal/domain/entities"

// OrderRepository определяет интерфейс для работы с заказами в базе данных
type OrderRepository interface {
	// Save сохраняет заказ в базу данных
	Save(order *entities.Order) error

	// GetByID получает заказ по ID
	GetByID(orderUID string) (*entities.Order, error)

	// GetAll получает все заказы
	GetAll() ([]*entities.Order, error)

	// Close закрывает соединение с базой данных
	Close() error
}
