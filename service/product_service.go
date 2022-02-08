package service

import (
	"net/http"
)

type ProductService interface {
	AddProductStockHandlerService(w http.ResponseWriter, r *http.Request)
	GetAvailableProductsHandlerService(w http.ResponseWriter, r *http.Request)
	OrderProductHandlerService(w http.ResponseWriter, r *http.Request)
}

type ProductInsightService interface {
	BestSellingProducts(w http.ResponseWriter, r *http.Request)
}
