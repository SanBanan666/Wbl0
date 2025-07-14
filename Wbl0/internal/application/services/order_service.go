package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"WbServis/Wbl0/internal/application/interfaces"
	"WbServis/Wbl0/internal/domain/entities"
)

type orderService struct {
	repository interfaces.OrderRepository
	cache      map[string]*entities.Order
	mutex      sync.RWMutex
}

func NewOrderService(repository interfaces.OrderRepository) interfaces.OrderService {
	return &orderService{
		repository: repository,
		cache:      make(map[string]*entities.Order),
	}
}

func (s *orderService) ProcessOrder(order *entities.Order) error {
	err := s.repository.Save(order)
	if err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	s.mutex.Lock()
	s.cache[order.OrderUID] = order
	s.mutex.Unlock()

	log.Printf("Order %s processed successfully", order.OrderUID)
	return nil
}

func (s *orderService) GetOrderByID(orderUID string) (*entities.Order, error) {
	s.mutex.RLock()
	if order, exists := s.cache[orderUID]; exists {
		s.mutex.RUnlock()
		log.Printf("Order %s found in cache", orderUID)
		return order, nil
	}
	s.mutex.RUnlock()

	order, err := s.repository.GetByID(orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from database: %w", err)
	}

	s.mutex.Lock()
	s.cache[orderUID] = order
	s.mutex.Unlock()

	log.Printf("Order %s loaded from database and cached", orderUID)
	return order, nil
}

func (s *orderService) RestoreCache() error {
	log.Println("Restoring cache from database...")

	orders, err := s.repository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all orders: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache = make(map[string]*entities.Order)

	for _, order := range orders {
		s.cache[order.OrderUID] = order
	}

	log.Printf("Cache restored with %d orders", len(orders))
	return nil
}

func (s *orderService) Close() error {
	return s.repository.Close()
}

func (s *orderService) ProcessMessage(message []byte) error {
	var order entities.Order
	err := json.Unmarshal(message, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	if order.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}

	return s.ProcessOrder(&order)
}
