// File: photoImageEditor.go

package controllers

import (
	"bufio"
	"bytes"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/jpeg"
	"io"
	"os"
)

type PhotoImageEditor struct {
	Latitude  float64 // GPS info
	Longitude float64 // GPS info

	ResizeScale  uint   // resize scale
	ResizedImage []byte // resized photo image

	fileSize int64 // original image file
}

// Processing Photo Image
func (c *PhotoImageEditor) ProcessPhotoImage(src *os.File) {

	// save original image file
	fileSave := make(chan bool)
	go func() {
		dst, err := os.Create(imageFileTemp)
		c.errorHandler(err)
		defer dst.Close()

		filesize, err := io.Copy(dst, src)
		c.fileSize = filesize
		c.errorHandler(err)

		fileSave <- false
	}()
	<-fileSave

	// extract GPS infomation
	gpsChannel := make(chan bool)
	go func() {
		buf1, err := os.Open(imageFileTemp)
		c.errorHandler(err)
		defer buf1.Close()

		c.Latitude, c.Longitude = c.extractGPSinfo(buf1)
		gpsChannel <- false
	}()

	// resize image data
	resizeChannel := make(chan bool)
	go func() {
		buf2, err := os.Open(imageFileTemp)
		c.errorHandler(err)
		defer buf2.Close()

		c.ResizedImage = c.resizePhoto(buf2, c.ResizeScale)
		resizeChannel <- false
	}()

	<-gpsChannel
	<-resizeChannel
}

// error handler
func (c *PhotoImageEditor) errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

// resize photo image
func (c *PhotoImageEditor) resizePhoto(file *os.File, scale uint) []byte {
	var img image.Image
	var err error

	// decode file
	buf := bufio.NewReader(file)
	img, _, err = image.Decode(buf)
	c.errorHandler(err)

	// resize image
	reseizedimg := resize.Resize(scale, 0, img, resize.Lanczos3)

	// make bytes data
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, reseizedimg, nil)
	c.errorHandler(err)

	// byte data
	byteData := buffer.Bytes()

	return byteData
}

// extract GPS infomation
func (c *PhotoImageEditor) extractGPSinfo(file *os.File) (float64, float64) {
	var lat, lng float64
	var err error

	// decode exif file
	exifInfo, errors := exif.Decode(file)
	if errors != nil {
		println(errors)
		lat = 0.0
		lng = 0.0
	} else {
		// extract lat and lng value from exif file
		lat, lng, err = exifInfo.LatLong()
		c.errorHandler(err)
	}

	return lat, lng
}
