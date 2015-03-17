// File: user.go

package controllers

import (
	"AirportNavigator/app/models"
	_ "fmt"
	"github.com/revel/revel"
)

type Users struct {
	*revel.Controller
}

func (c Users) makeAllItemList() []models.Shopitem {
	// stack of shop item object
	var shopitems []models.Shopitem

	// data base controller
	db := new(DBController)
	defer db.Close()

	db.Init(shopitemDB)
	objs := db.All()

	if len(objs) != 0 {
		// unmarshal
		for _, obj := range objs {

			// create mapper object
			mapper := new(MapperOfShopItem)
			mapper.ShopItemMapping(obj)

			// mapping result
			shopitem := mapper.object

			// stacking
			shopitems = append(shopitems, shopitem)
		}
	}

	return shopitems
}

func (c Users) Main() revel.Result {

	return c.Render()
}

func (c Users) Index(username string) revel.Result {
	// open csv file
	csv := new(CsvFileReader)
	csv.Open(csvFilePath)

	// airport name list
	airportNameList := csv.EnName

	// list of the Shop items
	shopitemlist := c.makeAllItemList()

	return c.Render(username, airportNameList, shopitemlist)
}

func (c Users) New(username string, email string, retypeemail string) revel.Result {
	return c.Render()
}
