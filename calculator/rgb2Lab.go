package calculator

import (
	"PixCSim/models"
)

/*
RGB2Lab :RGB to Lab calcualter
*/
func RGB2Lab(eachPatch8bitData map[int][]uint8, cspace models.ColorSpaceEnum, wpoint *models.WhitePoint, gamma float64) map[int][]float64 {

	// White point conversion matrix
	matElm := ConversionMatrixInit(cspace, wpoint)

	// Lab conversion matrix
	x, y, z := LabConversionMatrixInit(cspace, wpoint)

	// return
	return XYZ2Lab(RGB2XYZ(eachPatch8bitData, gamma, matElm), x, y, z)
}
