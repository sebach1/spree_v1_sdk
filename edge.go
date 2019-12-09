package spree

type ProductsEdge struct {
	Products    []*Product `json:"products,omitempty"`
	Count       int        `json:"count"`
	TotalCount  int        `json:"total_count"`
	CurrentPage int        `json:"current_page"`
	PerPage     int        `json:"per_page"`
	Pages       int        `json:"pages"`
}
