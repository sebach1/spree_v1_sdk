package spree

type Webhook struct {
	Site     string
	Resource *Resource
	Parent   *Resource
}

type Resource struct {
	Id   int
	Name string
}
