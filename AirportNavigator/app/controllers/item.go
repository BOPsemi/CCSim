// File:item.go

package controllers

import (
	"AirportNavigator/app/models"
	"bytes"
	_ "fmt"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"image/jpeg"
	"os"
	"strings"
)

type Items struct {
	*revel.Controller
}

// fetch object
func (c Items) fetchObjetct(info map[string]string, key bson.M) []interface{} {
	// open data base
	db := new(DBController)
	defer db.Close()

	db.Init(info)
	objs := db.Find(key)

	return objs
}

// setup image
func (c Items) setupImage(obj models.ShopitemImage) string {
	file := imageFileURL + obj.Uuid + ".jpg"

	// setup output stream
	outfile, _ := os.Create(file)
	defer outfile.Close()

	// make buffer
	buffer := bytes.NewBuffer(obj.Data)

	// make jpeg file
	img, _ := jpeg.Decode(buffer)
	err := jpeg.Encode(outfile, img, nil)
	if err != nil {
		panic(err)
	}

	println(file)

	return file
}

// make shop item array
func (c Items) mekeShopitemArray(itemlist string) []string {
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

// make shop item list
func (c Items) makeShopitemList(objs []string) string {
	var listStr string

	if len(objs) != 0 {
		count := len(objs) - 1
		for index, str := range objs {
			if index < count {
				listStr += str + ";" + "\n"
			} else {
				listStr += str + ";"
			}
		}
	} else {
		listStr = ""
	}

	return listStr
}

// ------------------
// Main Actios
// ------------------

// Add action
func (c Items) Add(
	airport int,
	name string,
	terminal string,
	floor int,
	comment string,
	imageFile *os.File,
	itemlist string) revel.Result {

	var imageFileUpload bool // check flag of image file

	var latitude float64  // GPS latitude
	var longitude float64 // GPS logitude
	var imageData []byte  // resized image data

	// flag set to false
	imageFileUpload = false

	if imageFile != nil {
		// create photo image editor object
		photoEditor := new(PhotoImageEditor)
		photoEditor.ResizeScale = resizeScaleFactor

		// processing image file
		photoEditor.ProcessPhotoImage(imageFile)

		// resize image
		imageData = photoEditor.ResizedImage

		// GPS infomation
		latitude = photoEditor.Latitude
		longitude = photoEditor.Longitude

		// set flat true
		imageFileUpload = true

	} else {

		// GPS infomation
		latitude = 0
		longitude = 0
	}

	// make information
	generator := new(ShopItemGenerator)
	generator.Init(airport, name, terminal, floor, comment, itemlist, latitude, longitude, imageFileUpload)

	// Airport Shop Item Object
	obj := generator.ShopItem

	// ------- Data base controller -------

	if imageData != nil {
		//  create object
		imageGenerator := new(ShopItemImageGenerator)
		imageGenerator.Init(imageData, obj.Uuid)

		imageObj := imageGenerator.ShopItemImage

		db := new(DBController)
		defer db.Close()

		db.Init(shopimageDB)
		db.Insert(imageObj)
	}

	// data base controller
	db := new(DBController)
	defer db.Close()

	db.Init(shopitemDB)
	db.Insert(obj)

	//println(obj.AirportInfo.AirportCode.IATA)

	return c.Redirect("/api/v1/users")
}

// Edit action
func (c Items) Edit(uuid string) revel.Result {

	// data
	var shopitems []models.Shopitem
	var shopimages []models.ShopitemImage

	// fetch key
	key := bson.M{"uuid": uuid}

	// fetch base data
	channelData := make(chan bool)
	go func() {
		objs := c.fetchObjetct(shopitemDB, key)
		mapper := new(ShopItemGenerator)

		for _, obj := range objs {
			mapper.ObjectMapping(obj)
			entity := mapper.ShopItem

			// stacking
			shopitems = append(shopitems, entity)
		}

		channelData <- false
	}()

	// fetch image data
	channelImage := make(chan bool)
	go func() {
		objs := c.fetchObjetct(shopimageDB, key)
		mapper := new(ShopItemImageGenerator)

		for _, obj := range objs {
			mapper.ObjectMapping(obj)
			entity := mapper.ShopItemImage

			// stacking
			shopimages = append(shopimages, entity)
		}

		channelImage <- false
	}()

	<-channelData
	<-channelImage

	// summarize fetch action
	var shopitem models.Shopitem
	var shopimage models.ShopitemImage
	var filePath string

	if len(shopitems) != 0 {
		// make jpg file
		shopimage = shopimages[0]
		filePath = c.setupImage(shopimage)
	} else {
		filePath = ""
	}

	shopitem = shopitems[0]
	itemlist := c.makeShopitemList(shopitem.Items)

	return c.Render(shopitem, itemlist, filePath)
}

// Update action
func (c Items) Update(
	name string,
	terminal string,
	floor int,
	comment string,
	itemlist string,
	button string,
	uuid string) revel.Result {

	// crate search key
	key := bson.M{"uuid": uuid}

	// setup database handler for shopItem
	dbItem := new(DBController)
	defer dbItem.Close()

	dbItem.Init(shopitemDB)

	// setup database handler for shopItemImage
	dbImage := new(DBController)
	defer dbImage.Close()

	dbImage.Init(shopimageDB)

	if button == "update" {
		// update information
		itemarray := c.mekeShopitemArray(itemlist)

		// make update phrase
		phrase := bson.M{
			"shopname": name,
			"terminal": terminal,
			"floor":    floor,
			"comment":  comment,
			"items":    itemarray,
		}

		// make key for shop item array
		phraseArray := bson.M{
			"items": bson.M{"$each": itemarray},
		}

		// update
		dbItem.Update(key, phrase)

		// update array
		dbItem.UpdateArray(key, phraseArray)

	} else {
		// delete entity
		dbItem.Remove(key)
		dbImage.Remove(key)
	}

	return c.Redirect("/api/v1/users")
}
