// File: mapperOfShopImage.go

package controllers

import (
	"AirportNavigator/app/models"
	"encoding/json"
)

// for shop item image
type MapperOfShopImage struct {
	mapperOfObject MapperOfObject
	object         models.ShopitemImage
}

// mapping
func (c *MapperOfShopImage) ShopImageMapping(entity interface{}) {

	err := json.Unmarshal(c.mapperOfObject.encoder(entity), &c.object)
	c.mapperOfObject.errorHandler(err)
}
