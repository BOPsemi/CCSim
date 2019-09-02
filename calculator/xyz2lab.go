package calculator

import (
	"PixCSim/models"
	"math"
)

/*
LabConversionMatrixInit :initialize Lab Conversion matrix
*/
func LabConversionMatrixInit(cspace models.ColorSpaceEnum, wpoint *models.WhitePoint) (xxelm, yyelm, zzelm float64) {
	// inline macro
	elmCalc := func(wElm float64) float64 {
		elm := 1.0 / wElm
		if elm < 0.0 {
			return 0.0
		}
		return elm
	}

	var xelm, yelm, zelm float64
	switch cspace {
	case models.CIE:
		xelm = elmCalc(wpoint.CIE.WhiteElm[0])
		yelm = elmCalc(wpoint.CIE.WhiteElm[1])
		zelm = elmCalc(wpoint.CIE.WhiteElm[2])

	case models.NTSC:
		xelm = elmCalc(wpoint.NTSC.WhiteElm[0])
		yelm = elmCalc(wpoint.NTSC.WhiteElm[1])
		zelm = elmCalc(wpoint.NTSC.WhiteElm[2])

	default:
		xelm = elmCalc(wpoint.SRGB.WhiteElm[0])
		yelm = elmCalc(wpoint.SRGB.WhiteElm[1])
		zelm = elmCalc(wpoint.SRGB.WhiteElm[2])
	}

	return xelm, yelm, zelm

}

/*
XYZ2Lab :convert xyz to Lab
*/
func XYZ2Lab(eachPatchXYZdata map[int][]float64, xelm, yelm, zelm float64) map[int][]float64 {
	labResults := make(map[int][]float64, 0)
	d1o3 := 1.0 / 3.0

	// Lab calculation
	lab := func(x, y, z float64) []float64 {
		L := 116*math.Pow(y, d1o3) - 16.0
		a := 500 * (math.Pow(x, d1o3) - math.Pow(y, d1o3))
		b := 200 * (math.Pow(y, d1o3) - math.Pow(z, d1o3))

		return []float64{L, a, b}
	}

	for index, xyzData := range eachPatchXYZdata {
		xxelm := xyzData[0] * xelm
		yyelm := xyzData[1] * yelm
		zzelm := xyzData[2] * zelm

		labResults[index] = lab(xxelm, yyelm, zzelm)
	}

	return labResults
}
