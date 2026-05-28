package models

type OrderEvent struct {
	UserID      int64  `json:"user_id"`
	OrderID     int64  `json:"order_id"`
	OrderStatus string `json:"order_status"`
}
