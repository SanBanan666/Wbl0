package interfaces

import "WbServis/Wbl0/internal/domain/entities"

// OrderRepository определяет интерфейс для работы с заказами в базе данных
type OrderRepository interface {
	Save(order *entities.Order) error

	GetByID(orderUID string) (*entities.Order, error)

	GetAll() ([]*entities.Order, error)

	Close() error
}
