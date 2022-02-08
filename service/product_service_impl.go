package service

import (
	"encoding/json"
	"github.com/awesomeProject/adapter"
	"github.com/awesomeProject/controller"
	"log"
	"net/http"
	"sync"
)

type ProductServiceImpl struct {
	ProductController controller.ProductServiceController
}

type ProductInsightServiceImpl struct {
	ProductInsightController controller.ProductInsightController
}

var mutex sync.Mutex

func (productServiceImpl ProductServiceImpl) AddProductStockHandlerService(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()

	products, err := adapter.ToProductRequest(req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addProductResponse, err := productServiceImpl.ProductController.AddProductsController(&products)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mutex.Unlock()
		return
	}
	addedProductsBytes, err := json.Marshal(addProductResponse)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mutex.Unlock()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(addedProductsBytes)
	mutex.Unlock()

}

func (productServiceImpl ProductServiceImpl) GetAvailableProductsHandlerService(w http.ResponseWriter, r *http.Request) {

	availableProduct, err := productServiceImpl.ProductController.GetAvailableProducts()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	productsBytes, err := json.Marshal(availableProduct)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(productsBytes)

}

func (productServiceImpl ProductServiceImpl) OrderProductHandlerService(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	w.Header().Set("Content-Type", "application/json")

	productOrderRequest, err := adapter.ToProductOrderRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orderResponse, err := productServiceImpl.ProductController.OrderProducts(&productOrderRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		orderResponse.OrderedProducts = nil
		orderResponse.OrderId = 0
		orderResponse.UnAvailableProducts = nil
	}

	orderResponseBytes, err := json.Marshal(orderResponse)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(orderResponseBytes)
	mutex.Unlock()

}

func (productInsightServiceImpl ProductInsightServiceImpl) BestSellingProducts(w http.ResponseWriter, r *http.Request) {
	bestSellerProducts, err := productInsightServiceImpl.ProductInsightController.GetBestSellerProducts()

	bestSellerProductsBytes, err := json.Marshal(bestSellerProducts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bestSellerProductsBytes)
}
