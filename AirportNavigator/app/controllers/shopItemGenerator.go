// Filee: shopItemGenerator.go

package controllers

import (
	"AirportNavigator/app/models"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ShopItemGenerator struct {
	ShopItem models.Shopitem
}

func (c *ShopItemGenerator) ObjectMapping(obj interface{}) {
	mapper := new(MapperOfShopItem)
	mapper.ShopItemMapping(obj)

	c.ShopItem = mapper.object
}

func (c ShopItemGenerator) errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func (c ShopItemGenerator) makeAirportObject(airportIndex int) *models.Airport {
	// open csv file
	csv := new(CsvFileReader)
	csv.Open(csvFilePath)

	// etract data
	iata := csv.IATA[airportIndex]
	icao := csv.ICAO[airportIndex]
	jpname := csv.JpName[airportIndex]
	enname := csv.EnName[airportIndex]

	// make code object
	chanelCode := make(chan bool)
	code := new(models.Code)
	go func() {
		code.IATA = iata
		code.ICAO = icao

		chanelCode <- false
	}()

	// make name object
	chanelName := make(chan bool)
	name := new(models.Name)
	go func() {
		name.JpName = jpname
		name.EnName = enname

		chanelName <- false
	}()

	<-chanelCode
	<-chanelName

	airport := new(models.Airport)
	airport.AirportCode = code
	airport.AirportName = name
	airport.Country = "jp"

	return airport
}

func (c ShopItemGenerator) makeUUID() string {
	out, err := exec.Command("uuidgen").Output()
	c.errorHandler(err)

	str := out[:36]

	return string(str)
}

func (c ShopItemGenerator) mekeGPSObject(lat float64, lng float64) *models.GPSInfo {
	gps := new(models.GPSInfo)
	gps.Lat = lat
	gps.Lng = lng

	return gps
}

func (c ShopItemGenerator) mekeItemlist(itemlist string) []string {
	items := strings.Split(itemlist, ";")
	itemNumber := len(items) - 1

	var list []string

	for index, item := range items {
		words := []byte(item)

		if index != 0 {
			if len(words) > 2 {
				words = words[2:len(words)]
			}
		}

		str := string(words)
		list = append(list, str)
	}

	list = list[:itemNumber]

	return list
}

func (c ShopItemGenerator) makeUpdateTimeStamp() string {
	return time.Now().Format("2006/01/02-15:04:05")
}

func (c ShopItemGenerator) makeImageURL(uuid string) string {
	url := imageFileURL + uuid + ".jpg"

	return url
}

func (c *ShopItemGenerator) Init(
	airport int,
	name string,
	terminal string,
	floor int,
	comment string,
	itemlist string,
	lat float64,
	lng float64,
	imageFileUpload bool) {

	// make airport shop item object
	c.ShopItem.AirportInfo = c.makeAirportObject(airport)
	c.ShopItem.Shopname = name
	c.ShopItem.Terminal = terminal
	c.ShopItem.Floor = floor
	c.ShopItem.GPS = c.mekeGPSObject(lat, lng)
	c.ShopItem.Uuid = c.makeUUID()
	c.ShopItem.Comment = comment
	c.ShopItem.Update = c.makeUpdateTimeStamp()
	c.ShopItem.Items = c.mekeItemlist(itemlist)

	if c.ShopItem.Uuid != "" {
		if imageFileUpload {
			url := c.makeImageURL(c.ShopItem.Uuid)
			c.ShopItem.URLs = append(c.ShopItem.URLs, url)
		}
	}

	fmt.Printf("Airport => %v\n", c.ShopItem.AirportInfo.AirportCode)
	fmt.Printf("Airport => %v\n", c.ShopItem.AirportInfo.AirportName)
	fmt.Printf("Image => %v\n", c.ShopItem.GPS)
	fmt.Printf("obj => %v\n", c.ShopItem)
}
