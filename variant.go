package spree

import (
	"encoding/json"
)

type VariantId int

type Variant struct {
	Id   VariantId `json:"id"`
	Name string    `json:"name"`
	Sku  string    `json:"sku"`

	Price     string `json:"price"`
	CostPrice string `json:"cost_price"`

	Weight string `json:"weight"`
	Height string `json:"height"`
	Width  string `json:"width"`
	Depth  string `json:"depth"`

	IsMaster        bool           `json:"is_master"`
	Slug            string         `json:"slug"`
	Description     string         `json:"description"`
	TrackInventory  bool           `json:"track_inventory"`
	OptionValues    []*OptionValue `json:"option_values"`
	Images          []*Image       `json:"images"`
	DisplayPrice    string         `json:"display_price"`
	OptionsText     string         `json:"options_text"`
	InStock         bool           `json:"in_stock"`
	IsBackorderable bool           `json:"is_backorderable"`
	TotalOnHand     int            `json:"total_on_hand"`
	IsDestroyed     bool           `json:"is_destroyed"`
	StockItems      []*StockItem   `json:"stock_items"`
}

type OptionType struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Presentation string `json:"presentation"`
	Position     int    `json:"position"`
}

type OptionValue struct {
	Id                     int    `json:"id"`
	Name                   string `json:"name"`
	Presentation           string `json:"presentation"`
	OptionTypeName         string `json:"option_type_name"`
	OptionTypeID           int    `json:"option_type_id"`
	OptionTypePresentation string `json:"option_type_presentation"`
}

// TODO: variants-only operations

func (v *Variant) StockItemByLocation(stockLocationId StockLocationId) (*StockItem, error) {
	for _, si := range v.StockItems {
		if si.StockLocationId == stockLocationId {
			return si, nil
		}
	}
	return nil, errStockItemNotFound
}

func (s *Spree) GetVariant(id VariantId) (*Variant, error) {
	params, err := s.paramsWithToken()
	if err != nil {
		return nil, err
	}
	URL, err := s.RouteTo("/variants/%v", params, id)
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
	v := &Variant{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
