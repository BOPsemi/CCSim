// File:shopitem.go

package models

// Airport international code
type Code struct {
	IATA string `bson:"iata" json:"iata"`
	ICAO string `bson:"icao" json:"icao"`
}

// Airport Name tags
type Name struct {
	EnName string `bson:"enname" json:"enname"`
	JpName string `bson:"jpname" json:"jpname"`
}

// Airport
type Airport struct {
	AirportCode *Code  `bson:"airportcode" json:"airportcode"`
	AirportName *Name  `bson:"airportname" json:"airportname"`
	Country     string `bson:"country" json:"country"`
}

// Location
type GPSInfo struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

// Shop Item
type Shopitem struct {
	// belonging Airport infomation
	AirportInfo *Airport `bson:"airportinfo" json:"airportinfo"`
	Shopname    string   `bson:"shopname" json:"shopname"`
	Terminal    string   `bson:"terminal" json:"terminal"`
	Floor       int      `bson:"floor" json:"floor"`
	GPS         *GPSInfo `bson:"gps" json:"gps"`
	// Item Infomation
	Uuid    string   `bson:"uuid" json:"uuid"`
	Comment string   `bson:"comment" json:"comment"`
	Update  string   `bson:"update" json:"update"`
	URLs    []string `bson:"urls" json:"urls"`
	Items   []string `bson:"items" json:"items"`
}
