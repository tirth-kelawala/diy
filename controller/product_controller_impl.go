package controller

import (
	"context"
	"database/sql"
	"errors"
	"github.com/awesomeProject/dao"
	"github.com/awesomeProject/model"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type ProductServiceControllerImpl struct {
	SqlDb   *sql.DB
	RedisDb *redis.Client
}

type ProductInsightControllerImpl struct {
	SqlDb *sql.DB
}

func (productController ProductServiceControllerImpl) AddProductsController(products *model.AddProductRequest) (addProductResponse model.AddProductResponse, err error) {
	addProductResponse.AddedProducts = make(map[string]float64)
	tx, err := productController.SqlDb.BeginTx(context.Background(), nil)
	if err != nil {
		return addProductResponse, err
	}

	defer tx.Rollback()

	for _, product := range products.Products {
		if product.Price > 0 && product.Quantity > 0 {
			err = dao.CreateProduct(&product, tx)
			if err != nil {
				return addProductResponse, err
			}
			quantity := dao.GetQuantityByNameAndPrice(&product, tx)

			if quantity > 0 {
				quantityInFloat := product.Quantity + quantity

				err = dao.UpdateQuantityForExistingStock(&quantityInFloat, time.Now(), &product.Name, tx)
				if err != nil {
					log.Panicln(err)
					return addProductResponse, err
				}
			} else {
				err = dao.InsertProductStock(&product, tx)
				if err != nil {
					log.Panicln(err)
					return addProductResponse, err
				}
			}

			var insertedProduct model.ProductDetail
			product.Price = 0

			err = dao.GetRedisObject(productController.RedisDb, product.Name, &insertedProduct)
			if err != nil {
				err = dao.SetRedisObject(productController.RedisDb, context.Background(), product.Name, product, 0)
			} else {
				product.Quantity = product.Quantity + insertedProduct.Quantity
				err = dao.SetRedisObject(productController.RedisDb, context.Background(), product.Name, product, 0)
			}
			if err != nil {
				log.Panicln(err)
				return addProductResponse, err
			}
			addProductResponse.AddedProducts[product.Name] = product.Quantity
		}
	}
	if err = tx.Commit(); err != nil {
		return addProductResponse, err
	}
	return addProductResponse, nil
}

func (productController ProductServiceControllerImpl) GetAvailableProducts() (productsResponse model.ProductsResponse, err error) {
	productsResponse.Products, err = dao.GetAvailableProducts(productController.RedisDb)
	if err != nil {
		return model.ProductsResponse{}, err
	}

	return productsResponse, nil
}

func (productController ProductServiceControllerImpl) OrderProducts(request *model.OrderProductsRequest) (orderProductsResponse model.OrderProductsResponse, err error) {

	tx, err := productController.SqlDb.BeginTx(context.Background(), nil)

	defer tx.Rollback()

	var isInsufficientQuantityOrdered bool

	//check if quantity is available
	for key, value := range request.ProductsOrder {
		if value > 0 {
			product, err := dao.GetStockById(key, productController.RedisDb)
			if err != nil || product.Quantity < value {
				isInsufficientQuantityOrdered = true
				orderProductsResponse.UnAvailableProducts = append(orderProductsResponse.UnAvailableProducts, key)
			}
		} else {
			isInsufficientQuantityOrdered = true
		}
	}

	if err != nil {
		return model.OrderProductsResponse{}, err
	}

	if !isInsufficientQuantityOrdered {

		orderId, err := dao.CreateOrderEntry(request, tx)

		orderProductsResponse.OrderId = orderId

		if err != nil {
			return model.OrderProductsResponse{}, err
		}

		for name, quantity := range request.ProductsOrder {
			var product model.ProductDetail

			err = dao.GetRedisObject(productController.RedisDb, name, &product)

			if err != nil {
				return model.OrderProductsResponse{}, err
			}

			if product.Quantity < quantity {
				err = errors.New("insufficient quantity")
				return model.OrderProductsResponse{}, err
			}

			product.Quantity = product.Quantity - quantity

			if product.Quantity == 0 {
				err = dao.DeleteByName(product.Name, productController.RedisDb, tx)
				if err != nil {
					return model.OrderProductsResponse{}, err
				}
			} else {
				productStockEntities, err := dao.GetProductStockByNameOrderByPrice(product.Name, tx)
				if err != nil {
					return model.OrderProductsResponse{}, err
				}

				for _, productStockEntity := range productStockEntities {
					if quantity <= 0 {
						break
					}
					if productStockEntity.Quantity <= quantity {
						quantity = quantity - productStockEntity.Quantity

						err := dao.DeleteByStockId(productStockEntity.StockId, tx)

						if err != nil {
							return model.OrderProductsResponse{}, err
						}
					} else {
						err = dao.UpdateProductStockByStockId(productStockEntity.Quantity-quantity, productStockEntity.StockId, tx)

						if err != nil {
							return model.OrderProductsResponse{}, err
						}
						quantity = 0
					}
				}
				err = dao.SetRedisObject(productController.RedisDb, context.Background(), product.Name, product, 0)
				if err != nil {
					return model.OrderProductsResponse{}, err
				}
			}
			orderProductsResponse.OrderedProducts = append(orderProductsResponse.OrderedProducts, name)
		}
	}

	err = tx.Commit()

	if err != nil {
		return model.OrderProductsResponse{}, err
	}

	return orderProductsResponse, nil
}

func (productInsightController ProductInsightControllerImpl) GetBestSellerProducts() (productsResponse model.ProductsResponse, err error) {
	topProducts, err := dao.GetTopOrderedProducts(productInsightController.SqlDb)
	productsResponse.Products = topProducts
	return
}
