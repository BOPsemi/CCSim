package calculator

import (
	"math"
)

/*
WHITEDIGITALNUMBER : patch 19 digital number
*/
const (
	WHITEDIGITALNUMBER = 243
)

// reverse gamma correction
func reverseGammaCorrection(rgbData []uint8, gamma float64) []float64 {

	// inline macro
	uint8Tofloat64 := func(digit uint8) float64 {
		return float64(digit) / WHITEDIGITALNUMBER
	}

	// inline macro

	deGamma := func(data float64) float64 {
		return math.Pow(data, gamma)
	}

	/*
		deGamma2 := func(data float64) float64 {
			if data > 0.04045 {
				c := (data + 0.055) / 1.055
				return math.Pow(c, gamma)
			}
			return data / 12.92
		}
	*/

	//
	r := deGamma(uint8Tofloat64(rgbData[0]))
	g := deGamma(uint8Tofloat64(rgbData[1]))
	b := deGamma(uint8Tofloat64(rgbData[2]))

	// retun r,g,b
	return []float64{r, g, b}
}

/*
DeGammaCorrection :de-gamma correction
*/
func DeGammaCorrection(eachPatchRGBData map[int][]uint8, gamma float64) map[int][]float64 {
	degammaData := make(map[int][]float64, 0)

	for index, rgbData := range eachPatchRGBData {
		degammaData[index] = reverseGammaCorrection(rgbData, gamma)
	}

	return degammaData
}
