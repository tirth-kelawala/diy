package adapter

import (
	"bytes"
	"encoding/json"
	"github.com/awesomeProject/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	INPUT = "http://stockoutmock.dunzo.in/api/v1"
)

func TestToProductOrderRequest(t *testing.T) {
	body, _ := json.Marshal(model.OrderProductsRequest{})
	req := httptest.NewRequest(http.MethodGet, INPUT, bytes.NewBuffer(body))

	_, err := ToProductOrderRequest(req)

	if err != nil {
		t.Fatal(err)
	}

}

func TestToProductRequest(t *testing.T) {
	body, _ := json.Marshal(model.AddProductRequest{})
	req := httptest.NewRequest(http.MethodGet, INPUT, bytes.NewBuffer(body))

	_, err := ToProductRequest(req)

	if err != nil {
		t.Fatal(err)
	}
}
