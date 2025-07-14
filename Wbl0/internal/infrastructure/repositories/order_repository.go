package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"WbServis/Wbl0/internal/application/interfaces"
	"WbServis/Wbl0/internal/domain/entities"
)

// OrderRepository реализует интерфейс для работы с заказами в PostgreSQL
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository создает новый экземпляр OrderRepository
func NewOrderRepository(db *sql.DB) interfaces.OrderRepository {
	return &OrderRepository{db: db}
}

// Save сохраняет заказ в базу данных
func (r *OrderRepository) Save(order *entities.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Сохраняем основной заказ
	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number = EXCLUDED.track_number,
			entry = EXCLUDED.entry,
			locale = EXCLUDED.locale,
			internal_signature = EXCLUDED.internal_signature,
			customer_id = EXCLUDED.customer_id,
			delivery_service = EXCLUDED.delivery_service,
			shard_key = EXCLUDED.shard_key,
			sm_id = EXCLUDED.sm_id,
			date_created = EXCLUDED.date_created,
			oof_shard = EXCLUDED.oof_shard,
			updated_at = CURRENT_TIMESTAMP
	`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Сохраняем информацию о доставке
	_, err = tx.Exec(`
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email,
			updated_at = CURRENT_TIMESTAMP
	`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	// Сохраняем информацию об оплате
	_, err = tx.Exec(`
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction,
			request_id = EXCLUDED.request_id,
			currency = EXCLUDED.currency,
			provider = EXCLUDED.provider,
			amount = EXCLUDED.amount,
			payment_dt = EXCLUDED.payment_dt,
			bank = EXCLUDED.bank,
			delivery_cost = EXCLUDED.delivery_cost,
			goods_total = EXCLUDED.goods_total,
			custom_fee = EXCLUDED.custom_fee,
			updated_at = CURRENT_TIMESTAMP
	`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	// Удаляем старые товары и добавляем новые
	_, err = tx.Exec("DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to delete old items: %w", err)
	}

	// Сохраняем товары
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name,
				sale, size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	return tx.Commit()
}

// GetByID получает заказ по ID
func (r *OrderRepository) GetByID(orderUID string) (*entities.Order, error) {
	// Получаем основной заказ
	var order entities.Order
	err := r.db.QueryRow(`
		SELECT order_uid, track_number, entry, locale, internal_signature,
			   customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
	`, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Получаем информацию о доставке
	err = r.db.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_uid = $1
	`, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	// Получаем информацию об оплате
	err = r.db.QueryRow(`
		SELECT transaction, request_id, currency, provider, amount,
			   payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments WHERE order_uid = $1
	`, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Получаем товары
	rows, err := r.db.Query(`
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1 ORDER BY id
	`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

// GetAll получает все заказы
func (r *OrderRepository) GetAll() ([]*entities.Order, error) {
	rows, err := r.db.Query(`
		SELECT order_uid FROM orders ORDER BY date_created DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get order UIDs: %w", err)
	}
	defer rows.Close()

	var orders []*entities.Order
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, fmt.Errorf("failed to scan order UID: %w", err)
		}

		order, err := r.GetByID(orderUID)
		if err != nil {
			log.Printf("Warning: failed to get order %s: %v", orderUID, err)
			continue
		}
		if order != nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// Close закрывает соединение с базой данных
func (r *OrderRepository) Close() error {
	return r.db.Close()
}
