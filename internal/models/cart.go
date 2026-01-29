package models

import "time"

type CartResponseItem struct {
	ProductID   int64   `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type CartResponse struct {
	Items     []CartResponseItem `json:"items"`
	Total     float64            `json:"total"`
	ItemCount int                `json:"item_count"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}
