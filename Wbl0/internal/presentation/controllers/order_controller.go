package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"WbServis/Wbl0/internal/application/dto"
	"WbServis/Wbl0/internal/application/interfaces"
)

type OrderController struct {
	orderService interfaces.OrderService
}

func NewOrderController(orderService interfaces.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (c *OrderController) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/order/")
	if path == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := c.orderService.GetOrderByID(path)
	if err != nil {
		log.Printf("Failed to get order %s: %v", path, err)
		http.Error(w, fmt.Sprintf("Order not found: %s", path), http.StatusNotFound)
		return
	}

	response := dto.OrderResponse{
		Order: order,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Order %s retrieved successfully", path)
}

func (c *OrderController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"service": "order-service",
	}

	json.NewEncoder(w).Encode(response)
}
