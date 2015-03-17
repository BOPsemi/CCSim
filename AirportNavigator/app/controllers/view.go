// File: view.go

package controllers

import (
	"AirportNavigator/app/models"
	"github.com/revel/revel"
	"os"
)

type Views struct {
	*revel.Controller
}

// make shopitem object array
func (c Views) makeShopitemObject(entity interface{}) models.Shopitem {
	mapper := new(MapperOfShopItem)
	mapper.ShopItemMapping(entity)

	return mapper.object
}

// json
func (c Views) JSON() revel.Result {
	// model stocker
	var list []models.Shopitem

	// setup database controller
	db := new(DBController)
	defer db.Close()

	db.Init(shopitemDB)
	objs := db.All()

	for _, obj := range objs {
		entity := c.makeShopitemObject(obj)
		list = append(list, entity)
	}

	return c.RenderJson(list)
}

// image file
func (c Views) Image(file string) revel.Result {
	// image file path
	path := imageFileURL + file

	// open image file
	image, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	return c.RenderFile(image, revel.Inline)
}
