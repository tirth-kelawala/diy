package adapter

import (
	"encoding/json"
	"github.com/awesomeProject/model"
	"io/ioutil"
	"net/http"
)

func ToProductRequest(req *http.Request) (productsRequest model.AddProductRequest, err error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return model.AddProductRequest{}, err
	}

	err = json.Unmarshal(bodyBytes, &productsRequest)
	if err != nil {
		return model.AddProductRequest{}, err
	}
	return
}

func ToProductOrderRequest(req *http.Request) (productOrderRequest model.OrderProductsRequest, err error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return productOrderRequest, err
	}

	err = json.Unmarshal(bodyBytes, &productOrderRequest)
	if err != nil {
		return productOrderRequest, err
	}
	return productOrderRequest, nil
}
