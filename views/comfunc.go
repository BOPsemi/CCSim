package views

import (
	"errors"
	"image"
	"image/png"
	"os"
	"strings"
)

// check the last suffix
func checkLastWordInPath(str string) string {
	if !strings.HasSuffix(str, "/") {
		newStr := str + "/"
		return newStr
	}
	return str
}

/*
Save image.Image object as PNG
*/
func saveImageAsPNG(filepath string, filename string, img image.Image) (bool, error) {
	// check input
	if filepath == "" || filename == "" {
		return false, errors.New("save file path is empty")

	}
	if img == nil {
		return false, errors.New("image file body is empty")
	}

	// open file
	openFilePath := checkLastWordInPath(filepath) + filename + ".png"
	file, err := os.OpenFile(openFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		return false, err
	}

	// save PNG file
	err = png.Encode(file, img)
	if err != nil {
		return false, err
	}

	return true, nil
}
