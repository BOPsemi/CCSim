package controllers

import (
	"PixCSim/iotool"
	"PixCSim/models"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

// get file list in folder
func getFileListIn(folder string, ext string) []string {
	var listupFileExt string

	if folder == "" {
		return nil
	}

	if ext == "" {
		listupFileExt = "*.csv"
	} else {
		listupFileExt = ext
	}

	// get working dir
	workDir := os.Getenv("APPATH")

	// make loading file path from workDir
	loadingFileDir := workDir + folder

	// get files in load dir
	getFileList := func(dir string) []string {
		list, err := filepath.Glob(dir + listupFileExt)
		if err != nil {
			log.Println("Work file reading error")
			return nil
		}
		return list
	}

	return getFileList(loadingFileDir)
}

// read CSV file
func readCSVFile(filepath string, c chan [][]string) {
	// variables
	var err error
	var file *os.File
	var data = make([][]string, 0)

	// initialize io handler
	ioHandler := iotool.NewIOHandler()

	// open file
	file, err = ioHandler.OpenFile(filepath)
	defer file.Close()
	if err != nil {
		c <- data
	}

	// read csv file
	data, err = ioHandler.ReadCSV(file)
	if err != nil {
		c <- data
	}

	// return read data
	c <- data
}

// serialize calculation result
func serializeColorPatchCode(data map[int]models.ColorPatch) map[int][]uint8 {
	if data == nil {
		return nil
	}

	converted := make(map[int][]uint8, 0)
	for index, rgbData := range data {
		converted[index] = []uint8{rgbData.Channel.R, rgbData.Channel.G, rgbData.Channel.B}
	}

	return converted
}

func createDummyImage(data []uint8) image.Image {
	context := gg.NewContext(1, 1)
	context.DrawRectangle(0.0, 0.0, 1.0, 1.0)
	context.SetRGBA255(int(data[0]), int(data[1]), int(data[2]), 255)
	context.Fill()

	return context.Image()
}

func extract8bitCodeFrom(img *image.RGBA) []uint8 {
	r, g, b, _ := img.At(0, 0).RGBA()

	return []uint8{uint8(r), uint8(g), uint8(b)}
}

func deMapLinearMatElm(data map[string]float64) []float64 {
	stocker := make([]float64, 6)
	for index, element := range data {
		switch index {
		case "a":
			stocker[0] = element
		case "b":
			stocker[1] = element
		case "c":
			stocker[2] = element
		case "d":
			stocker[3] = element
		case "e":
			stocker[4] = element
		case "f":
			stocker[5] = element
		}
	}

	return stocker
}

func mapLinearMatElm(data []float64) map[string]float64 {
	stocker := make(map[string]float64, 0)
	for index, element := range data {
		switch index {
		case 0:
			stocker["a"] = element
		case 1:
			stocker["b"] = element
		case 2:
			stocker["c"] = element
		case 3:
			stocker["d"] = element
		case 4:
			stocker["e"] = element
		case 5:
			stocker["f"] = element
		}
	}

	return stocker
}
