// File:init.go

package controllers

// Constant
const (
	csvFilePath   = "AirportNavigator/resources/AirportList.csv"
	imageFileTemp = "AirportNavigator/tmp/buffer.jpg"
	imageFileURL  = "AirportNavigator/tmp/"
)

const (
	resizeScaleFactor = 680
)

// Index for CSV file Open
type Index int

func (i Index) String() int {
	switch i {
	case IATA:
		return 0
	case ICAO:
		return 1
	case JPname:
		return 2
	case ENname:
		return 3
	}

	return 999
}

const (
	IATA Index = iota
	ICAO
	JPname
	ENname
)

// --- DB parameters ---
// DEFINE: DBparameters for shop

const (
	DB_SKIP  int = 0
	DB_LIMIT int = 10
)
const (
	DB_URL_SHOP        string = "localhost"
	DB_SHOP            string = "airportshop"
	DB_COLLECTION_SHOP string = "airportshop"
)

// DEFINE: DBparemeters for shop image
const (
	DB_URL_SHOPIMAGE        string = "localhost"
	DB_SHOPIMAGE            string = "airportshopimage"
	DB_COLLECTION_SHOPIMAGE string = "airportshopimage"
)
const (
	DB_URL_SHOPITEM        string = "localhost"
	DB_SHOPITEM            string = "airportshopitem"
	DB_COLLECTION_SHOPITEM string = "airportshopitem"
)

// --- make DB map ---
var (
	// DEFINE: DB map shop model
	shopDB = map[string]string{
		"url":        DB_URL_SHOP,
		"name":       DB_SHOP,
		"collection": DB_COLLECTION_SHOP,
	}
	// DEFINE: DB map for shop image
	shopimageDB = map[string]string{
		"url":        DB_URL_SHOPIMAGE,
		"name":       DB_SHOPIMAGE,
		"collection": DB_COLLECTION_SHOPIMAGE,
	}
	// DEFINE: DB map for shop item
	shopitemDB = map[string]string{
		"url":        DB_URL_SHOPITEM,
		"name":       DB_SHOPITEM,
		"collection": DB_COLLECTION_SHOPITEM,
	}
)
