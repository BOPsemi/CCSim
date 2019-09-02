package calculator

import "math"

/*
GammaCorrection :calculate gamma
*/
func GammaCorrection(eachPatchRGBData map[int][]float64, gamma float64) map[int][]float64 {
	results := make(map[int][]float64, 0)
	invGamma := 1.0 / gamma

	// normarize signal level
	normRef := eachPatchRGBData[19][1]

	for index, rgbData := range eachPatchRGBData {
		r := math.Pow(rgbData[0]/normRef, invGamma)
		g := math.Pow(rgbData[1]/normRef, invGamma)
		b := math.Pow(rgbData[2]/normRef, invGamma)

		results[index] = []float64{r, g, b}
	}

	return results
}
