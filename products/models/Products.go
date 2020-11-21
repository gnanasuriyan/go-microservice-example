package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	SKU         string    `json:"sku"`
	CreatedOn   time.Time `json:"-"`
	UpdatedOn   time.Time `json:"-"`
	DeletedOn   time.Time `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

type ProductList []*Product

func (pl *ProductList) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(pl)
}

var productList = ProductList{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC(),
		UpdatedOn:   time.Now().UTC(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "def1234",
		CreatedOn:   time.Now().UTC(),
		UpdatedOn:   time.Now().UTC(),
	},
}

var nextProductId = 3

func GetProductList() ProductList {
	return productList
}

func AddProduct(product *Product) ProductList {
	product.ID = nextProductId
	nextProductId++
	productList = append(productList, product)
	return productList
}

func findProductIndex(productId int) (int, error) {
	for index, product := range productList {
		if product.ID == productId {
			return index, nil
		}
	}
	return -1, errors.New("product doesn't exist")
}
func UpdateProduct(IDProduct int, product *Product) error {
	index, err := findProductIndex(IDProduct)
	if err != nil {
		return err
	}
	product.ID = IDProduct
	productList[index] = product
	return nil
}
