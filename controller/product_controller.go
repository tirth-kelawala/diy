package controller

import "github.com/awesomeProject/model"

type ProductServiceController interface {
	AddProductsController(products *model.AddProductRequest) (addProductResponse model.AddProductResponse, err error)
	GetAvailableProducts() (productsResponse model.ProductsResponse, err error)
	OrderProducts(request *model.OrderProductsRequest) (orderProductsResponse model.OrderProductsResponse, err error)
}

type ProductInsightController interface {
	GetBestSellerProducts() (productsResponse model.ProductsResponse, err error)
}
