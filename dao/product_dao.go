package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/awesomeProject/model"
	"github.com/go-redis/redis"
	"github.com/kisielk/sqlstruct"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func CreateProduct(product *model.ProductDetail, tx *sql.Tx) (err error) {
	sqlStatement := "INSERT INTO product (name, description) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING;"

	_, err = tx.Exec(sqlStatement, product.Name, product.Description)
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

func GetQuantityByNameAndPrice(product *model.ProductDetail, tx *sql.Tx) (quantity float64) {
	sqlStatement := "SELECT quantity FROM product_stock WHERE name = $1 AND price = $2"
	tx.QueryRow(sqlStatement, product.Name, product.Price).Scan(&quantity)
	return quantity
}

func UpdateQuantityForExistingStock(quantity *float64, timestamp time.Time, name *string, tx *sql.Tx) (err error) {
	sqlStatement := "UPDATE product_stock SET quantity = $1, updated_at = $2 WHERE name = $3"
	_, err = tx.Exec(sqlStatement, quantity, timestamp, name)
	if err != nil {
		return err
	}
	return nil
}

func InsertProductStock(product *model.ProductDetail, tx *sql.Tx) (err error) {
	sqlStatement := `
		INSERT INTO product_stock (name, quantity, price, updated_at)
		VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(sqlStatement, product.Name, product.Quantity, product.Price, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func GetAvailableProducts(redisDb *redis.Client) ([]model.ProductDetail, error) {
	var availableProducts []model.ProductDetail

	iter := redisDb.Scan(context.Background(), 0, "*", 0).Iterator()

	for iter.Next(context.Background()) {
		var product model.ProductDetail
		err := GetRedisObject(redisDb, iter.Val(), &product)
		if err != nil {
			return availableProducts, err
		}
		availableProducts = append(availableProducts, product)
	}
	return availableProducts, nil
}

func GetRedisObject(redisDb *redis.Client, key string, dest interface{}) (err error) {
	val, err := redisDb.Get(context.Background(), key).Result()

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	return err
}

func SetRedisObject(redisDb *redis.Client, ctx context.Context, key string, object interface{}, expiration time.Duration) (err error) {
	objectBytes, err := json.Marshal(object)

	err = redisDb.Set(ctx, key, objectBytes, expiration).Err()
	return err
}

func DeleteByName(name string, redisDb *redis.Client, tx *sql.Tx) (err error) {
	sqlStatement := "DELETE FROM product_stock WHERE name = $1"
	_, err = tx.Exec(sqlStatement, name)
	if err != nil {
		return err
	}

	err = redisDb.Del(context.Background(), name).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetProductStockByNameOrderByPrice(name string, tx *sql.Tx) (productStockEntities []model.ProductStockEntity, err error) {
	sqlStatement := "SELECT * FROM product_stock where name=$1 ORDER BY price"

	products, err := tx.Query(sqlStatement, name)
	if err != nil {
		return nil, err
	}

	for products.Next() {
		var productStockEntity model.ProductStockEntity

		sqlstruct.Scan(&productStockEntity, products)

		if err != nil {
			return nil, err
		}

		productStockEntities = append(productStockEntities, productStockEntity)
	}
	err = products.Close()
	return productStockEntities, err
}

func UpdateProductStockByStockId(quantity float64, stockId int64, tx *sql.Tx) (err error) {
	sqlStatement := "UPDATE product_stock SET quantity = $1, updated_at = $2 WHERE stock_id = $3"
	product, err := tx.Query(sqlStatement, quantity, time.Now(), stockId)
	if err != nil {
		return err
	}
	err = product.Close()
	return err
}

func DeleteByStockId(stockId int64, tx *sql.Tx) (err error) {
	sqlStatement := "DELETE FROM product_stock where stock_id=$1"
	_, err = tx.Exec(sqlStatement, stockId)
	if err != nil {
		return err
	}
	return err
}

func GetStockById(name string, redisDb *redis.Client) (product model.ProductDetail, err error) {
	err = GetRedisObject(redisDb, name, &product)
	if err != nil {
		return model.ProductDetail{}, err
	}
	return product, nil
}

func CreateOrderEntry(request *model.OrderProductsRequest, tx *sql.Tx) (orderId int64, err error) {
	sqlStatement := "INSERT INTO product_order (updated_at) VALUES ($1) RETURNING order_id"

	err = tx.QueryRow(sqlStatement, time.Now()).Scan(&orderId)

	if err != nil {
		return -1, err
	}

	for key, value := range request.ProductsOrder {
		sqlStatement = "INSERT INTO ordered_product VALUES ($1, $2, $3)"
		_, err = tx.Exec(sqlStatement, orderId, key, value)
		if err != nil {
			return -1, err
		}
	}

	return orderId, err
}

func GetTopOrderedProducts(sqlDb *sql.DB) (bestSellerProducts []model.ProductDetail, err error) {
	sqlStatement := "select name,SUM(quantity) as quantity from product_order as po inner join ordered_product as op on po.order_id=op.order_id where updated_at >= NOW() - INTERVAL '1 hour' group by name order by quantity desc limit 5"
	orderedProducts, err := sqlDb.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	for orderedProducts.Next() {
		var productStock model.ProductDetail

		sqlstruct.Scan(&productStock, orderedProducts)

		if err != nil {
			return nil, err
		}

		bestSellerProducts = append(bestSellerProducts, productStock)
	}
	err = orderedProducts.Close()

	return
}
