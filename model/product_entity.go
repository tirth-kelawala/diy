package model

import "time"

type ProductEntity struct {
	Name        string
	Description string
}

type ProductDetail struct {
	Name        string
	Description string  `json:",omitempty"`
	Price       float64 `json:",omitempty"`
	Quantity    float64
}

type AddProductRequest struct {
	Products []ProductDetail
	Comment  string
	Username string
}

type AddProductResponse struct {
	AddedProducts map[string]float64
}

type ProductsResponse struct {
	Products []ProductDetail
}

type OrderProductsRequest struct {
	ProductsOrder map[string]float64
}

type OrderProductsResponse struct {
	OrderId             int64    `json:",omitempty"`
	OrderedProducts     []string `json:",omitempty"`
	UnAvailableProducts []string `json:",omitempty"`
}

type OrderEntity struct {
	OrderId    int
	Created_at time.Time
}

type ProductStockEntity struct {
	StockId   int64     `sql:"stock_id"`
	Name      string    `sql:"name"`
	Price     float64   `sql:"price"`
	Quantity  float64   `sql:"quantity"`
	UpdatedAt time.Time `sql:"updated_at"`
}
