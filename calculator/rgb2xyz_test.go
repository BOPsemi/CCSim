package calculator

import (
	"PixCSim/models"
	"fmt"
	"testing"
)

func Test_ConversionMatrix(t *testing.T) {
	wpoint := models.NewWhitePoint(models.SRGB)
	matElm := ConversionMatrixInit(models.SRGB, wpoint)

	mocData := make(map[int][]uint8, 0)
	data := []uint8{100, 52, 52}
	mocData[1] = data
	mocData[2] = data

	//xyzResult := RGB2XYZ(mocData, 2.2, matElm)
	x, y, z := LabConversionMatrixInit(models.SRGB, wpoint)
	labResult := XYZ2Lab(RGB2XYZ(mocData, 2.2, matElm), x, y, z)

	fmt.Println(labResult)

}

func Test_RGB2Lab(t *testing.T) {
	wpoint := models.NewWhitePoint(models.SRGB)

	mocData := make(map[int][]uint8, 0)
	data := []uint8{100, 52, 52}
	mocData[1] = data
	mocData[2] = data

	// RGB -> Lab
	results := RGB2Lab(mocData, models.SRGB, wpoint, 2.2)
	fmt.Println(results)

}
