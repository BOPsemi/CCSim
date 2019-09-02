package controllers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSimulation(t *testing.T) {
	obj := NewSimulation()
	assert.NotNil(t, obj)

	linearMat := obj.LoadedLinearMatrix()

	linearMatResults := obj.ApplyLinearMatrix(linearMat)
	whiteBalancedResults, rgain, bgain := obj.AdjustWhiteBalance(linearMatResults, 22)
	gammaCorrectedResults := obj.CorrectGamma(whiteBalancedResults, 2.2)
	digitizedResults := obj.DigitizeRaw(gammaCorrectedResults)
	serializedResults := serializeColorPatchCode(digitizedResults.Data)

	// brightness
	//adjustedBrightness := obj.AdjustBrightness(serializedResults, 1.0)

	// saturation
	adjustedSaturation := obj.AdjustSaturation(serializedResults, -0.5)

	//

	for index, data := range serializedResults {
		//result := adjustedBrightness[index]
		result := adjustedSaturation[index]
		fmt.Println(data, result)
	}

	fmt.Println(rgain, bgain)

}
