// File: shopitemimage.go

package models

import ()

type ShopitemImage struct {
	Data   []byte `bson:"data" json:"data"`
	Uuid   string `bson:"uuid" json:"uuid"`
	Update string `bson:"update" json:"update"`
}
