package spree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type ProductId int

type Product struct {
	Id          ProductId `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`

	Price        string `json:"price,omitempty"`
	DisplayPrice string `json:"display_price,omitempty"`

	ShippingCategoryId ShippingCategoryId `json:"shipping_category_id,omitempty"`
	TaxonIds           []int              `json:"taxon_ids,omitempty"`

	TotalOnHand int `json:"total_on_hand,omitempty"`

	OptionTypes []*OptionType `json:"option_types,omitempty"`

	ProductProperties []*Property `json:"product_properties,omitempty"`

	Master   *Variant   `json:"master,omitempty"`
	Variants []*Variant `json:"variants,omitempty"`
	// MetaDescription    interface{} `json:"meta_description"`
	// MetaKeywords       interface{} `json:"meta_keywords"`

	Slug        string    `json:"slug,omitempty"`
	AvailableOn time.Time `json:"available_on,omitempty"`
}

func (p *Product) VariantsIncludingMaster() []*Variant {
	return append(p.Variants, p.Master)
}

func ParsePrice(strPrice string) (float64, error) {
	price, err := strconv.ParseFloat(strPrice, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (s *Spree) SetProduct(prod *Product) (newProd *Product, err error) {
	if prod == nil {
		return nil, errNilProduct
	}
	if prod.Id == 0 {
		newProd, err = s.createProduct(prod)
	} else {
		newProd, err = s.updateProduct(prod)
	}
	return
}

func (s *Spree) createProduct(prod *Product) (*Product, error) {
	params, err := s.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := s.RouteTo("/products", params)
	if err != nil {
		return nil, err
	}
	jsonProd, err := json.Marshal(prod)
	if err != nil {
		return nil, err
	}
	resp, err := s.Post(URL, bytes.NewReader(jsonProd))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newProd := &Product{}
	err = json.NewDecoder(resp.Body).Decode(newProd)
	if err != nil {
		return nil, err
	}
	return newProd, nil
}

func (s *Spree) updateProduct(prod *Product) (*Product, error) {
	params, err := s.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := s.RouteTo("/products/%v", params, prod.Id)

	if err != nil {
		return nil, err
	}
	prod.Id = 0 // Unset id since its in the route
	jsonProd, err := json.Marshal(prod)
	if err != nil {
		return nil, err
	}
	resp, err := s.Put(URL, bytes.NewReader(jsonProd))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newProd := &Product{}
	err = json.NewDecoder(resp.Body).Decode(newProd)
	if err != nil {
		return nil, err
	}
	return newProd, nil
}

func (s *Spree) GetProduct(id ProductId) (*Product, error) {
	params, err := s.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := s.RouteTo("/products/%v", params, id)
	if err != nil {
		return nil, err
	}
	resp, err := s.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	prod := &Product{}
	err = json.NewDecoder(resp.Body).Decode(prod)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *Spree) GetProducts(page int) (*ProductsEdge, error) {
	params, err := s.paramsWithToken()
	if err != nil {
		return nil, err
	}
	params.Set("page", string(page))
	URL, err := s.RouteTo("/products", params)
	if err != nil {
		return nil, err
	}
	resp, err := s.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	edge := &ProductsEdge{}
	err = json.NewDecoder(resp.Body).Decode(edge)
	if err != nil {
		return nil, err
	}
	return edge, nil
}

func (s *Spree) getProducts(page int, errCh chan<- error, prodCh chan<- []*Product) {
	params, err := s.paramsWithToken()
	if err != nil {
		errCh <- err
		return
	}
	params.Set("page", string(page))
	URL, err := s.RouteTo("/products", params)
	if err != nil {
		errCh <- err
		return
	}
	resp, err := s.Get(URL)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		errCh <- errFromReader(resp.Body)
		return
	}
	edge := &ProductsEdge{}
	err = json.NewDecoder(resp.Body).Decode(edge)
	if err != nil {
		errCh <- err
		return
	}
	prodCh <- edge.Products
}

func (s *Spree) FetchProducts() ([]*Product, error) {
	edge, err := s.GetProducts(1)
	if err != nil {
		return nil, err
	}
	prods, pages := edge.Products, edge.Pages

	prodCh := make(chan []*Product)
	errCh := make(chan error)

	for currentPage := 2; currentPage <= pages; currentPage++ {
		go s.getProducts(currentPage, errCh, prodCh)
	}

	for returnedReqs := pages - 1; returnedReqs > 0; returnedReqs-- {
		select {
		case prodsArr := <-prodCh:
			prods = append(prods, prodsArr...)
		case err := <-errCh:
			return nil, err
		}
	}
	return prods, nil
}

func (s *Spree) DeleteProduct(id ProductId) error {
	params, err := s.paramsWithToken()
	if err != nil {
		return err
	}
	URL, err := s.RouteTo("/products/%v", params, id)
	if err != nil {
		return err
	}
	resp, err := s.Delete(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return errFromReader(resp.Body)
	}
	return nil
}

type ShippingCategoryId int

// TODO: add images as needed
func NewProduct(name string, price float64, shippingCategoryId ShippingCategoryId) (*Product, error) {
	prod := &Product{Name: name, Price: TransformPrice(price), ShippingCategoryId: shippingCategoryId}
	err := prod.validate()
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func TransformPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)

}

func (p *Product) validate() error {
	price, err := strconv.ParseFloat(p.Price, 64)
	if err != nil {
		return err
	}
	if price == 0 {
		return errNilPrice
	}
	if p.Name == "" {
		return errNilProductName
	}
	if p.ShippingCategoryId == 0 {
		return errNilShippingCategoryId
	}
	return nil
}
