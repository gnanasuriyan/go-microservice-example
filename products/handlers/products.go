package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gnanasuriyan/go-micro-services-http/products/models"
)

type Products struct {
	logger *log.Logger
}

func NewProductHandler(logger *log.Logger) *Products {
	return &Products{logger: logger}
}
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.logger.Printf("http method: %s", r.Method)
	if r.Method == http.MethodGet {
		p.listProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.updateProduct(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) listProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("listing products")
	productList := models.GetProductList()
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("adding new product")
	newProduct := new(models.Product)
	if err := newProduct.FromJSON(r.Body); err != nil {
		p.logger.Printf("unable to read product information: %v", r.Body)
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
	p.logger.Printf("add product request payload: %#v", newProduct)
	productList := models.AddProduct(newProduct)
	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
	}
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("updating product", r.URL.Path)

	// get the id from URL path parameter
	regExp := regexp.MustCompile(`/([0-9]+)`)
	groups := regExp.FindAllStringSubmatch(r.URL.Path, -1)

	p.logger.Println("received URL parameters", groups)

	if len(groups) != 1 {
		p.logger.Printf("more than one id was received: %v", groups)
		http.Error(rw, "invalid product id", http.StatusBadRequest)
		return
	}

	if len(groups[0]) != 2 {
		p.logger.Printf("more than one captured group: %v", groups[0])
		http.Error(rw, "invalid product id", http.StatusBadRequest)
		return
	}
	productId, err := strconv.Atoi(groups[0][1])
	if err != nil {
		p.logger.Println("invalid product id", groups[0][1])
		http.Error(rw, "invalid product id", http.StatusBadRequest)
		return
	}

	updateProduct := new(models.Product)
	if err := updateProduct.FromJSON(r.Body); err != nil {
		p.logger.Printf("unable to read product information: %v", r.Body)
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
	p.logger.Printf("add product request payload: %#v", updateProduct)
	if err = models.UpdateProduct(productId, updateProduct); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
