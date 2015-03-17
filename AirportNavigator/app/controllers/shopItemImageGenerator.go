// File :shopItemImageGenerator.go

package controllers

import (
	"AirportNavigator/app/models"
	"time"
)

type ShopItemImageGenerator struct {
	ShopItemImage models.ShopitemImage
}

func (c *ShopItemImageGenerator) ObjectMapping(obj interface{}) {
	mapper := new(MapperOfShopImage)
	mapper.ShopImageMapping(obj)

	c.ShopItemImage = mapper.object
}

func (c ShopItemImageGenerator) errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func (c ShopItemImageGenerator) makeUpdateTimeStamp() string {
	return time.Now().Format("2006/01/02-15:04:05")
}

func (c *ShopItemImageGenerator) Init(imageData []byte, uuid string) {
	c.ShopItemImage.Data = imageData
	c.ShopItemImage.Uuid = uuid
	c.ShopItemImage.Update = c.makeUpdateTimeStamp()
}
