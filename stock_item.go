package spree

import (
	"bytes"
	"encoding/json"
)

type StockItemId int

type StockLocationId int

type StockItem struct {
	Id                StockItemId     `json:"id"`
	CountOnHand       int             `json:"count_on_hand"`
	StockLocationId   StockLocationId `json:"stock_location_id"`
	Backorderable     bool            `json:"backorderable"`
	Available         bool            `json:"available"`
	StockLocationName string          `json:"stock_location_name"`
	Force             bool            `json:"force,omitempty"`
}

func (s *Spree) UpdateStock(id VariantId, stockLocationId StockLocationId, stock int) (*StockItem, error) {
	v, err := s.GetVariant(id)
	if err != nil {
		return nil, err
	}
	si, err := v.StockItemByLocation(stockLocationId)
	if err != nil {
		return nil, err
	}
	params, err := s.paramsWithToken()
	if err != nil {
		return nil, err
	}
	si.Force = true // overrides actual stock; unless the api will perform an append
	si.CountOnHand = stock
	si.Id = 0 // unset id due its on param
	URL, err := s.RouteTo("/stock_locations/%v/stock_items/%v", params, stockLocationId, si.Id)
	if err != nil {
		return nil, err
	}
	jsonStockItem, err := json.Marshal(si)
	if err != nil {
		return nil, err
	}
	resp, err := s.Put(URL, bytes.NewReader(jsonStockItem))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errFromReader(resp.Body)
	}
	newStockItem := &StockItem{}
	err = json.NewDecoder(resp.Body).Decode(newStockItem)
	if err != nil {
		return nil, err
	}
	return newStockItem, nil
}
