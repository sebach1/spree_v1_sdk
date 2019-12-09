package spree

type PropertyId int
type Property struct {
	Id           PropertyId
	ProductId    ProductId
	PropertyId   int
	Value        string
	PropertyName string
}
