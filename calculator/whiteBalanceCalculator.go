package calculator

/*
WhiteBalanceCalculater :calculate white balance
*/
func WhiteBalanceCalculater(eachPatchRGBdata map[int][]float64, refPatchNumber int) (adjusted map[int][]float64, rgain, bgain float64) {
	result := make(map[int][]float64, 0)

	// extract reference patch
	refPatch := eachPatchRGBdata[refPatchNumber]

	// calculate gain
	gvalue := refPatch[1]
	redGain := refPatch[1] / refPatch[0]
	blueGain := refPatch[1] / refPatch[2]

	for index, rgbData := range eachPatchRGBdata {
		r := rgbData[0] * redGain / gvalue
		g := rgbData[1] * 1.0 / gvalue
		b := rgbData[2] * blueGain / gvalue

		result[index] = []float64{r, g, b}
	}
	return result, redGain, blueGain
}
