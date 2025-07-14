package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "orders"

	orders := []map[string]interface{}{
		{
			"order_uid":    "b563feb7b2b84b6test",
			"track_number": "WBILMTESTTRACK",
			"entry":        "WBIL",
			"delivery": map[string]interface{}{
				"name":    "Test Testov",
				"phone":   "+9720000000",
				"zip":     "2639809",
				"city":    "Kiryat Mozkin",
				"address": "Ploshad Mira 15",
				"region":  "Kraiot",
				"email":   "test@gmail.com",
			},
			"payment": map[string]interface{}{
				"transaction":   "b563feb7b2b84b6test",
				"request_id":    "",
				"currency":      "USD",
				"provider":      "wbpay",
				"amount":        1817,
				"payment_dt":    1637907727,
				"bank":          "alpha",
				"delivery_cost": 1500,
				"goods_total":   317,
				"custom_fee":    0,
			},
			"items": []map[string]interface{}{
				{
					"chrt_id":      9934930,
					"track_number": "WBILMTESTTRACK",
					"price":        453,
					"rid":          "ab4219087a764ae0btest",
					"name":         "Mascaras",
					"sale":         30,
					"size":         "0",
					"total_price":  317,
					"nm_id":        2389212,
					"brand":        "Vivienne Sabo",
					"status":       202,
				},
			},
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        "test",
			"delivery_service":   "meest",
			"shardkey":           "9",
			"sm_id":              99,
			"date_created":       "2021-11-26T06:22:19Z",
			"oof_shard":          "1",
		},
		{
			"order_uid":    "test-order-1",
			"track_number": "TRACK001",
			"entry":        "TEST",
			"delivery": map[string]interface{}{
				"name":    "John Doe",
				"phone":   "+1234567890",
				"zip":     "12345",
				"city":    "New York",
				"address": "123 Main St",
				"region":  "NY",
				"email":   "john@example.com",
			},
			"payment": map[string]interface{}{
				"transaction":   "txn-001",
				"request_id":    "req-001",
				"currency":      "USD",
				"provider":      "stripe",
				"amount":        2500,
				"payment_dt":    time.Now().Unix(),
				"bank":          "chase",
				"delivery_cost": 500,
				"goods_total":   2000,
				"custom_fee":    0,
			},
			"items": []map[string]interface{}{
				{
					"chrt_id":      12345,
					"track_number": "TRACK001",
					"price":        1000,
					"rid":          "rid-001",
					"name":         "Test Product 1",
					"sale":         0,
					"size":         "M",
					"total_price":  1000,
					"nm_id":        67890,
					"brand":        "Test Brand",
					"status":       202,
				},
				{
					"chrt_id":      12346,
					"track_number": "TRACK001",
					"price":        1000,
					"rid":          "rid-002",
					"name":         "Test Product 2",
					"sale":         0,
					"size":         "L",
					"total_price":  1000,
					"nm_id":        67891,
					"brand":        "Test Brand",
					"status":       202,
				},
			},
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        "customer-1",
			"delivery_service":   "fedex",
			"shardkey":           "1",
			"sm_id":              1,
			"date_created":       time.Now().Format(time.RFC3339),
			"oof_shard":          "1",
		},
		{
			"order_uid":    "test-order-2",
			"track_number": "TRACK002",
			"entry":        "TEST",
			"delivery": map[string]interface{}{
				"name":    "Jane Smith",
				"phone":   "+0987654321",
				"zip":     "54321",
				"city":    "Los Angeles",
				"address": "456 Oak Ave",
				"region":  "CA",
				"email":   "jane@example.com",
			},
			"payment": map[string]interface{}{
				"transaction":   "txn-002",
				"request_id":    "req-002",
				"currency":      "EUR",
				"provider":      "paypal",
				"amount":        1500,
				"payment_dt":    time.Now().Unix(),
				"bank":          "wells_fargo",
				"delivery_cost": 300,
				"goods_total":   1200,
				"custom_fee":    0,
			},
			"items": []map[string]interface{}{
				{
					"chrt_id":      23456,
					"track_number": "TRACK002",
					"price":        600,
					"rid":          "rid-003",
					"name":         "Premium Product",
					"sale":         20,
					"size":         "S",
					"total_price":  480,
					"nm_id":        78901,
					"brand":        "Premium Brand",
					"status":       202,
				},
				{
					"chrt_id":      23457,
					"track_number": "TRACK002",
					"price":        600,
					"rid":          "rid-004",
					"name":         "Premium Product 2",
					"sale":         20,
					"size":         "M",
					"total_price":  480,
					"nm_id":        78902,
					"brand":        "Premium Brand",
					"status":       202,
				},
			},
			"locale":             "en",
			"internal_signature": "",
			"customer_id":        "customer-2",
			"delivery_service":   "ups",
			"shardkey":           "2",
			"sm_id":              2,
			"date_created":       time.Now().Format(time.RFC3339),
			"oof_shard":          "2",
		},
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	for i, order := range orders {
		orderBytes, err := json.Marshal(order)
		if err != nil {
			log.Printf("Failed to marshal order %d: %v", i+1, err)
			continue
		}
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(order["order_uid"].(string)),
			Value: sarama.ByteEncoder(orderBytes),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf(" Ошибка отправки заказа %v: %v", order["order_uid"], err)
		} else {
			fmt.Printf(" Заказ %v отправлен (partition=%d, offset=%d)\n", order["order_uid"], partition, offset)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n Все тестовые заказы отправлены!")
	os.Exit(0)
}
