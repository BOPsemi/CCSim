package controllers

import (
	"PixCSim/models"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mocLinearMat() map[string]float64 {
	moc := make(map[string]float64, 0)
	moc["a"] = 0.220
	moc["b"] = 0.040
	moc["c"] = 0.550
	moc["d"] = 0.440
	moc["e"] = 0.015
	moc["f"] = 0.350

	return moc
}

func Test_NewOptimizer(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	matElm := mocLinearMat()
	optimizer := NewOptimizer(matElm, 2.2, models.SRGB)
	assert.NotNil(t, optimizer)

	optimizer.Run(2.0)
	ccData, linerMat, redGain, blueGain, hueAngle, brightness, saturation, deltaE := optimizer.OptimizedResults()
	fmt.Println(ccData)
	fmt.Println(linerMat)
	fmt.Println(redGain, blueGain)
	fmt.Println(hueAngle)
	fmt.Println(brightness, saturation)
	fmt.Println(deltaE)

	eachPatchDeltaE := optimizer.OptimizedEachPatchDeltaE()
	for index := 1; index <= len(eachPatchDeltaE); index++ {
		fmt.Println(index, eachPatchDeltaE[index])
	}
}

func Test_deMap(t *testing.T) {
	matElm := make(map[string]float64, 0)
	matElm["a"] = 100.0
	matElm["b"] = 200.0
	matElm["c"] = 300.0
	matElm["d"] = 400.0
	matElm["e"] = 500.0
	matElm["f"] = 600.0

	elmMat := deMapLinearMatElm(matElm)
	fmt.Println(elmMat)

	assert.NotEqual(t, 0, len(elmMat))
}

func Test_map(t *testing.T) {
	mocData := []float64{100, 200, 300, 400, 500, 600}

	elmMat := mapLinearMatElm(mocData)
	fmt.Println(elmMat)

	assert.Equal(t, 6, len(elmMat))
}
