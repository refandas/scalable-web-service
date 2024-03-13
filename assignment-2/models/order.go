package models

import (
	"fmt"
	"github.com/refandas/scalable-web-service/assignment-2/database"
	"time"
)

type CreateOrderInput struct {
	OrderedAt    time.Time `json:"orderedAt" validate:"required"`
	CustomerName string    `json:"customerName" validate:"required,max=50"`
	Items        []struct {
		ItemCode    string `json:"itemCode" validate:"required,max=10"`
		Description string `json:"description" validate:"required,max=50"`
		Quantity    int    `json:"quantity" validate:"required,min=1"`
	} `json:"items" validate:"required,dive"`
}

type UpdateOrderInput struct {
	OrderedAt    time.Time `json:"orderedAt" validate:"omitempty"`
	CustomerName string    `json:"customerName" validate:"omitempty,max=50"`
	Items        []struct {
		ItemCode    string `json:"itemCode" validate:"required,max=10"`
		Description string `json:"description" validate:"required,max=50"`
		Quantity    int    `json:"quantity" validate:"required,min=1"`
	} `json:"items" validate:"omitempty,dive"`
}

type Order struct {
	ID           int64     `json:"id"`
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items"`
}

func (order *CreateOrderInput) CreateOrder() error {
	db := database.GetDB()
	defer db.Close()

	var orderId int64
	err := db.QueryRow(`
		INSERT INTO orders (ordered_at, customer_name) 
		VALUES ($1, $2) 
		RETURNING id
	`, time.Now(), order.CustomerName).Scan(&orderId)

	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := db.Exec(`
			INSERT INTO items (code, description, quantity, order_id) 
			VALUES ($1, $2, $3, $4)
		`, item.ItemCode, item.Description, item.Quantity, orderId)

		if err != nil {
			return err
		}
	}
	return nil
}

func (order *CreateOrderInput) ToOrderResponse(orderId int64) Order {
	items := make([]Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = Item{
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
		}
	}
	return Order{
		ID:           orderId, // Assuming ID is not available in CreateOrderInput
		OrderedAt:    order.OrderedAt,
		CustomerName: order.CustomerName,
		Items:        items,
	}
}

func (order *UpdateOrderInput) UpdateOrder(orderId int64) (*Order, error) {
	db := database.GetDB()
	defer db.Close()

	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM orders WHERE id = $1)", orderId).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("order with ID %d not found", orderId)
	}

	tx, err := db.Begin()
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = tx.Exec(`
		UPDATE orders SET ordered_at = $1, customer_name = $2
		WHERE id = $3
	`, order.OrderedAt, order.CustomerName, orderId)

	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM items WHERE order_id = $1", orderId)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return nil, err
		}
		return nil, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO items (order_id, code, description, quantity)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return nil, err
		}
		return nil, err
	}
	defer stmt.Close()

	for _, item := range order.Items {
		_, err = stmt.Exec(orderId, item.ItemCode, item.Description, item.Quantity)
		if err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				return nil, err
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	orderResult := Order{ID: orderId}
	updatedOrder, err := orderResult.GetOrder()
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (order *UpdateOrderInput) ToOrderResponse(orderId int64) Order {
	items := make([]Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = Item{
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
		}
	}
	return Order{
		ID:           orderId,
		OrderedAt:    order.OrderedAt,
		CustomerName: order.CustomerName,
		Items:        items,
	}
}

func (order *Order) GetOrder() (*Order, error) {
	db := database.GetDB()
	defer db.Close()

	err := db.QueryRow(`
		SELECT id, customer_name, ordered_at 
		FROM orders 
		WHERE id = $1
	`, order.ID).Scan(&order.ID, &order.CustomerName, &order.OrderedAt)

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT code, description, quantity FROM items WHERE order_id = $1", order.ID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ItemCode, &item.Description, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items

	return order, nil
}

func GetOrders() ([]Order, error) {
	db := database.GetDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, customer_name, ordered_at FROM orders")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.OrderedAt); err != nil {
			return nil, err
		}

		itemRows, err := db.Query(`
			SELECT code, description, quantity 
			FROM items 
			WHERE order_id = $1
		`, order.ID)

		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		var items []Item
		for itemRows.Next() {
			var item Item
			if err := itemRows.Scan(&item.ItemCode, &item.Description, &item.Quantity); err != nil {
				return nil, err
			}
			items = append(items, item)
		}

		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func DeleteOrder(orderId int64) error {
	db := database.GetDB()

	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM orders WHERE id = $1)", orderId).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("order with ID %d not found", orderId)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM items WHERE order_id = $1", orderId)

	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return err
		}
		return err
	}

	_, err = tx.Exec("DELETE FROM orders WHERE id = $1", orderId)

	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return errTx
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
