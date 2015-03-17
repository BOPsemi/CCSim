// File: jsonToObjectMapper.go

package controllers

import (
	"AirportNavigator/app/models"
	"encoding/json"
)

// For Shop Item
type MapperOfShopItem struct {
	mapperOfObject MapperOfObject
	object         models.Shopitem
}

// mapping
func (c *MapperOfShopItem) ShopItemMapping(entity interface{}) {

	err := json.Unmarshal(c.mapperOfObject.encoder(entity), &c.object)
	c.mapperOfObject.errorHandler(err)
}
