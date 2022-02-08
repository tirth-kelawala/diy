package dao

import (
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/awesomeProject/model"
	"github.com/go-redis/redismock"
	"regexp"
	"testing"
	"time"
)

func TestCreateProduct(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO product").WithArgs("p1", "d1").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := sqlDb.BeginTx(context.TODO(), nil)
	defer sqlDb.Close()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := model.ProductDetail{
		Name:        "p1",
		Description: "d1",
	}

	err = CreateProduct(&product, tx)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
}

func TestGetQuantityByNameAndPrice(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := model.ProductDetail{
		Name:        "p1",
		Description: "d1",
		Quantity:    10,
		Price:       10,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT quantity FROM product_stock").WithArgs(product.Name, product.Price).WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(product.Quantity))
	tx, err := sqlDb.BeginTx(context.TODO(), nil)
	defer sqlDb.Close()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	quan := GetQuantityByNameAndPrice(&product, tx)

	if quan != product.Quantity {
		t.Fatalf("quantity '%s'", err)
	}
}

func TestUpdateQuantityForExistingStock(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := model.ProductDetail{
		Name:        "p1",
		Description: "d1",
		Quantity:    10,
		Price:       10,
	}

	currentTime := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE product_stock").WithArgs(product.Quantity, currentTime, product.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	tx, err := sqlDb.BeginTx(context.TODO(), nil)
	defer sqlDb.Close()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = UpdateQuantityForExistingStock(&product.Quantity, currentTime, &product.Name, tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetProductStockByNameOrderByPrice(t *testing.T) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := model.ProductStockEntity{
		StockId:   1,
		Name:      "test",
		Price:     1,
		Quantity:  1,
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM product_stock where name=$1 ORDER BY price")).WithArgs(product.Name).WillReturnRows(sqlmock.NewRows([]string{"stock_id", "name", "price", "quantity", "updated_at"}).AddRow(product.StockId, product.Name, product.Price, product.Quantity, product.UpdatedAt))
	tx, err := sqlDb.BeginTx(context.TODO(), nil)
	defer sqlDb.Close()

	if err != nil {
		t.Fatal(err)
	}

	_, err = GetProductStockByNameOrderByPrice(product.Name, tx)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteByName(t *testing.T) {
	redisClient, mockClient := redismock.NewClientMock()
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	product := model.ProductDetail{
		Name:        "p1",
		Description: "d1",
		Quantity:    10,
		Price:       10,
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM product_stock WHERE name=&1").WithArgs(product.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	tx, err := sqlDb.BeginTx(context.TODO(), nil)
	defer sqlDb.Close()

	if err != nil {
		t.Fatal(err)
	}

	mockClient.ExpectDel(product.Name).SetVal(1)

	err = DeleteByName(product.Name, redisClient, tx)

	if err != nil {
		t.Fatal(err)
	}

}

func TestGetRedisObject(t *testing.T) {
	redisClient, mockClient := redismock.NewClientMock()

	var testKey = "key"
	var testValue = model.ProductDetail{
		Name:        "",
		Description: "",
		Price:       10,
		Quantity:    10,
	}

	b, err := json.Marshal(testValue)

	mockClient.ExpectGet(testKey).SetVal(string(b))

	var retrievedValue model.ProductDetail

	err = GetRedisObject(redisClient, testKey, &retrievedValue)
	if err != nil || retrievedValue.Price != testValue.Price {
		t.Fatal(err)
	}

	err = GetRedisObject(redisClient, "", &retrievedValue)
	if err == nil {
		t.Fatal(err)
	}
}

func TestSetRedisObject(t *testing.T) {
	redisClient, mockClient := redismock.NewClientMock()
	var testKey = "key"
	var testValue = model.ProductDetail{
		Name:        "",
		Description: "",
		Price:       10,
		Quantity:    10,
	}

	objectBytes, err := json.Marshal(testValue)

	if err != nil {
		t.Fatal(err)
	}

	mockClient.ExpectSet(testKey, objectBytes, 0).SetVal("")

	err = SetRedisObject(redisClient, context.TODO(), testKey, testValue, 0)

	if err != nil {
		t.Fatal(err)
	}
}
