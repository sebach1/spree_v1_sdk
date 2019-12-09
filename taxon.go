package spree

type Classification struct {
	TaxonID  int `json:"taxon_id"`
	Position int `json:"position"`
	Taxon    *Taxon
}

type Taxon struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	PrettyName string   `json:"pretty_name"`
	Permalink  string   `json:"permalink"`
	ParentId   int      `json:"parent_id"`
	TaxonomyId int      `json:"taxonomy_id"`
	Taxons     []*Taxon `json:"taxons"`
}
