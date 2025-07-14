package dto

import "WbServis/Wbl0/internal/domain/entities"

// OrderResponse представляет ответ API для заказа
type OrderResponse struct {
	Order *entities.Order `json:"order"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// OrderRequest представляет запрос для создания заказа
type OrderRequest struct {
	Order *entities.Order `json:"order"`
}
