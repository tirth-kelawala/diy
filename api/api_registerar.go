package api

import (
	"github.com/awesomeProject/factory"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func HandleRequests() {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/product-management/products", factory.ProductService.AddProductStockHandlerService).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/product-management/available-products", factory.ProductService.GetAvailableProductsHandlerService).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/product-management/order", factory.ProductService.OrderProductHandlerService).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/product-management/products/best-seller", factory.ProductInsightService.BestSellingProducts).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}
