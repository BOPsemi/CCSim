package controllers

import (
	"PixCSim/calculator"
	"PixCSim/models"

	"github.com/anthonynsimon/bild/adjust"

	"gonum.org/v1/gonum/mat"
)

/*
COLORCHARTNUM :total number of color chart
*/
const (
	COLORCHARTNUM = 24
)

/*
Simulation :interface of simulation
*/
type Simulation interface {
	// simulation
	ApplyLinearMatrix(linearMat map[string]float64) map[int][]float64
	AdjustWhiteBalance(eachPatchRGBdata map[int][]float64, refPatchNumber int) (adjusted map[int][]float64, redGain, blueGain float64)
	CorrectGamma(eachPatchRGBdata map[int][]float64, gamma float64) map[int][]float64
	DigitizeRaw(rawEachPatchRGBdata map[int][]float64) *models.MacbethColorCode
	AdjustBrightness(eachPatchRGB8bitData map[int][]uint8, brightness float64) map[int][]uint8
	AdjustSaturation(eachPatchRGB8bitData map[int][]uint8, saturation float64) map[int][]uint8
	AdjustHue(eachPatchRGB8bitData map[int][]uint8, angle int) map[int][]uint8

	// getter
	LoadedLinearMatrix() map[string]float64
	StandardMacbethColorCode() map[int]models.ColorPatch
}

// structure
type simulation struct {
	startupHandler StartUp // startup handler
	setupHandler   Setup   // setup handler

	colorChartReflection     map[int][]float64
	deviceWavelengthResponse map[int][]float64

	/*
		this object defines the object before integration
		wavelegth, [r, gr, gb, b]
	*/
	signal struct {
		r  map[int][]float64
		gr map[int][]float64
		gb map[int][]float64
		b  map[int][]float64
	}

	/*
		This object defines channel response after integration
		[patch01, patch02, ...]
	*/
	eachPatchChannelResponse struct {
		r  []float64
		gr []float64
		gb []float64
		b  []float64
	}

	/*
		The loaded linear matrix
	*/
	loadedLinearMat map[string]float64
}

/*
NewSimulation :initializer of simulation
*/
func NewSimulation() Simulation {
	// initialize object
	obj := new(simulation)

	// call initializating function
	if !obj.initSimulation() {
		return nil
	}

	// return initialized object
	return obj
}

/*
ApplyLinearMatrix :apply linear matrix, use the inputted linear matrix
*/
func (sim *simulation) ApplyLinearMatrix(linearMat map[string]float64) map[int][]float64 {
	return calculator.LinearMatixCalculator(
		sim.eachPatchChannelResponse.r,
		sim.eachPatchChannelResponse.gr,
		sim.eachPatchChannelResponse.gb,
		sim.eachPatchChannelResponse.b,
		linearMat,
	)
}

/*
AdjustWhiteBalance :adjust white balance automatically
*/
func (sim *simulation) AdjustWhiteBalance(eachPatchRGBdata map[int][]float64, refPatchNumber int) (adjusted map[int][]float64, redGain, blueGain float64) {
	return calculator.WhiteBalanceCalculater(
		eachPatchRGBdata,
		refPatchNumber,
	)
}

/*
CorrectGamma : call gamma correction function
*/
func (sim *simulation) CorrectGamma(eachPatchRGBdata map[int][]float64, gamma float64) map[int][]float64 {
	return calculator.GammaCorrection(
		eachPatchRGBdata,
		gamma,
	)
}

/*
DigitizeRaw :digitize simulation data
*/
func (sim *simulation) DigitizeRaw(rawEachPatchRGBdata map[int][]float64) *models.MacbethColorCode {
	return models.NewMacbethColorCodeFromRaw(rawEachPatchRGBdata)
}

/*
AdjustBrightness :adjust image brightness
*/
func (sim *simulation) AdjustBrightness(eachPatchRGB8bitData map[int][]uint8, brightness float64) map[int][]uint8 {
	adjustedData := make(map[int][]uint8, 0)

	for index, rgbData := range eachPatchRGB8bitData {
		img := createDummyImage(rgbData)
		adjustedImg := adjust.Brightness(img, brightness)

		adjustedData[index] = extract8bitCodeFrom(adjustedImg)
	}

	return adjustedData
}

/*
AdjustSaturation :adjust image saturation
*/
func (sim *simulation) AdjustSaturation(eachPatchRGB8bitData map[int][]uint8, saturation float64) map[int][]uint8 {
	adjustedData := make(map[int][]uint8, 0)

	for index, rgbData := range eachPatchRGB8bitData {
		img := createDummyImage(rgbData)
		adjustedImg := adjust.Saturation(img, saturation)

		adjustedData[index] = extract8bitCodeFrom(adjustedImg)
	}

	return adjustedData
}

/*
AdjustHue :adjust Hue
*/
func (sim *simulation) AdjustHue(eachPatchRGB8bitData map[int][]uint8, angle int) map[int][]uint8 {
	adjustedData := make(map[int][]uint8, 0)

	for index, rgbData := range eachPatchRGB8bitData {
		img := createDummyImage(rgbData)
		adjustedImg := adjust.Hue(img, angle)

		adjustedData[index] = extract8bitCodeFrom(adjustedImg)
	}

	return adjustedData
}

/*
LoadedLinearMatrix : return loaded linear matrix
*/
func (sim *simulation) LoadedLinearMatrix() map[string]float64 {
	return sim.loadedLinearMat
}

func (sim *simulation) StandardMacbethColorCode() map[int]models.ColorPatch {
	return sim.startupHandler.StandardMacbethColorCode()
}

// channel signal
func (sim *simulation) channelSignal() (R, Gr, Gb, B []float64) {

	sumOfEachPatch := func(channel map[int][]float64) []float64 {
		sum := mat.NewDense(1, COLORCHARTNUM, nil)
		// scan wavelength
		for _, patchData := range channel {
			a := mat.NewDense(1, COLORCHARTNUM, patchData)
			sum.Add(sum, a)
		}

		return sum.RawRowView(0)
	}

	r := sumOfEachPatch(sim.signal.r)
	gr := sumOfEachPatch(sim.signal.gr)
	gb := sumOfEachPatch(sim.signal.gb)
	b := sumOfEachPatch(sim.signal.b)

	return r, gr, gb, b
}

// calculate each patch response
func (sim *simulation) eachPatchSignal() (R, Gr, Gb, B map[int][]float64) {
	r := make(map[int][]float64, 0)
	gr := make(map[int][]float64, 0)
	gb := make(map[int][]float64, 0)
	b := make(map[int][]float64, 0)

	for wavelength, devRes := range sim.deviceWavelengthResponse {
		// define stocker, the each patch signals are accumulated in slice
		var rstocker []float64
		var grstocker []float64
		var gbstocker []float64
		var bstocker []float64

		// scan each patch reflection, and calculate device signal
		for _, reflection := range sim.colorChartReflection[wavelength] {
			rstocker = append(rstocker, devRes[0]*reflection)
			grstocker = append(grstocker, devRes[1]*reflection)
			gbstocker = append(gbstocker, devRes[2]*reflection)
			bstocker = append(bstocker, devRes[3]*reflection)
		}

		if rstocker != nil && len(rstocker) != 0 {
			r[wavelength] = rstocker
			gr[wavelength] = grstocker
			gb[wavelength] = gbstocker
			b[wavelength] = bstocker
		}

	}

	return r, gr, gb, b
}

/*
InitSimulation :initialize simulation object
*/
func (sim *simulation) initSimulation() bool {
	var err error

	// initialize startup and setup
	sim.startupHandler = NewStartUp()
	sim.setupHandler = NewSetup()

	// get data
	sim.colorChartReflection = make(map[int][]float64, 0)
	sim.colorChartReflection, err = sim.startupHandler.ColorChartReflection()
	if err != nil {
		sim.colorChartReflection = nil
		return false
	}
	sim.deviceWavelengthResponse = make(map[int][]float64, 0)
	sim.deviceWavelengthResponse, err = sim.setupHandler.DeviceWavelengthResponse()
	if err != nil {
		sim.deviceWavelengthResponse = nil
		return false
	}
	sim.loadedLinearMat = sim.setupHandler.LinearMatrixEelemtns()

	// calculate signal from device
	r, gr, gb, b := sim.eachPatchSignal()
	sim.signal.r = r
	sim.signal.gr = gr
	sim.signal.gb = gb
	sim.signal.b = b

	sim.eachPatchChannelResponse.r, sim.eachPatchChannelResponse.gr, sim.eachPatchChannelResponse.gb, sim.eachPatchChannelResponse.b = sim.channelSignal()

	return true
}
