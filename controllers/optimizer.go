package controllers

import (
	"PixCSim/calculator"
	"PixCSim/models"
	"errors"
	"log"
	"math"
	"math/rand"
	"time"
)

/*
Optimizer :define optimizer interface
*/
type Optimizer interface {
	// setter
	SetOptimizerCondition(bachSize, maxTrial int, rangePercent float64) (bool, error)
	SetOptimizerCoeffecients(alpha, beta1, beta2, epsilon float64) (bool, error)

	// getter
	OptimizedResults() (ccData map[int][]uint8, linearMat map[string]float64, redGain, blueGain float64, hueAngle int, brightness, saturation, deltaEAve float64)
	OptimizedEachPatchDeltaE() map[int]float64

	// run optimizer
	Run(targetDeltaE float64) bool
}

// struct
type optimizer struct {
	// conditions
	condition struct {
		bachSize     int     // 1 back size
		maxTrial     int     // Maximum trial number
		rangePercent float64 // range value for calculating div, uniti is percent
		hue          struct {
			minAngle int // Hue rotation min angle
			maxAngle int // Hue rotation max angle
		}
		brightness struct {
			low  float64 // brightness sweep low value
			high float64 // brightness sweep high value
		}
		saturation struct {
			low  float64 // saturation sweep low value
			high float64 // saturation sweep high value
		}
	}

	// Adam inital condition
	adamCondition struct {
		alpha   float64
		beta1   float64
		beta2   float64
		epsilon float64
	}

	// k-values for deltaE calculation
	kvals []float64

	// referece color chart data
	refCCData  map[int][]uint8
	refLabData map[int][]float64

	// result color chart data
	resultCCData     map[int][]uint8
	resultDeltaE     float64
	resultBrightness float64
	resultSaturation float64
	resultLinarMat   map[string]float64
	resultHueAngle   int
	resultEachDeltaE map[int]float64

	redGain  float64
	blueGain float64

	// stocker of initial value
	originalData     map[int][]uint8
	initLinearMatElm []float64
	gamma            float64
	colorSpace       models.ColorSpaceEnum
	whitePoint       *models.WhitePoint

	// simulation
	simulation Simulation
}

/*
NewOptimizer : initializer of object
*/
func NewOptimizer(initLinearMatElm map[string]float64, gamma float64, cspace models.ColorSpaceEnum) Optimizer {
	obj := new(optimizer)

	// initialize properties
	obj.initLinearMatElm = deMapLinearMatElm(initLinearMatElm)
	obj.gamma = gamma
	obj.colorSpace = cspace
	obj.whitePoint = models.NewWhitePoint(cspace)

	// initialize conditions
	obj.condition.bachSize = 10
	obj.condition.maxTrial = 50
	obj.condition.rangePercent = 0.03

	obj.condition.hue.minAngle = -5
	obj.condition.hue.maxAngle = 5

	obj.condition.brightness.low = -0.5
	obj.condition.brightness.high = 0.5

	obj.condition.saturation.low = -0.5
	obj.condition.saturation.high = 0.5

	// Adam initial condition
	obj.adamCondition.alpha = 0.01
	obj.adamCondition.beta1 = 0.9
	obj.adamCondition.beta2 = 0.999
	obj.adamCondition.epsilon = 1.0e-3

	// kvalues for deltaE calculation
	obj.kvals = []float64{1.0, 1.0, 1.0}

	// initialize results
	obj.resultCCData = make(map[int][]uint8, 0)
	obj.resultLinarMat = make(map[string]float64, 0)

	// initialize object
	obj.simulation = NewSimulation()

	// initialize reference CCData
	obj.refCCData = serializeColorPatchCode(obj.simulation.StandardMacbethColorCode())
	obj.refLabData = calculator.RGB2Lab(obj.refCCData, cspace, obj.whitePoint, gamma)

	return obj
}

// simulate Image
func (op *optimizer) simulateImage(linearMat map[string]float64, gamma float64) (data map[int][]uint8, rgain, bgain float64) {
	/*
		Simulate image from inputted linear matrix
	*/
	linearMatResults := op.simulation.ApplyLinearMatrix(linearMat)
	whiteBalancedResults, redGain, blueGain := op.simulation.AdjustWhiteBalance(linearMatResults, 22)
	gammaCorrectedResults := op.simulation.CorrectGamma(whiteBalancedResults, gamma)
	digitizedResults := op.simulation.DigitizeRaw(gammaCorrectedResults)
	serializedResults := serializeColorPatchCode(digitizedResults.Data)

	return serializedResults, redGain, blueGain
}

// gradient calculator
func (op *optimizer) gradient(elm []float64, elmIndex int, target float64) float64 {
	// initialize buffer
	elmSize := len(elm)
	minusElm := make([]float64, elmSize)
	plusElm := make([]float64, elmSize)

	// copy elm silice
	copy(minusElm, elm)
	copy(plusElm, elm)

	// define inline
	// --- Calculate minum linear matrix
	minus := func(elm []float64, index int) {
		ans := elm[index] - elm[index]*op.condition.rangePercent
		if ans < 0.0 {
			ans = 0.0
		}
		elm[index] = ans
	}

	// --- Calculate plus linear matrix
	plus := func(elm []float64, index int) {
		ans := elm[index] + elm[index]*op.condition.rangePercent
		if ans > 1.0 {
			ans = 1.0
		}
		elm[index] = ans
	}

	// --- Simulate image and return Lab
	simulateImageFromMat := func(c chan map[int][]float64, linearMat []float64) {
		ccData, _, _ := op.simulateImage(mapLinearMatElm(linearMat), op.gamma)
		lab := calculator.RGB2Lab(ccData, op.colorSpace, op.whitePoint, op.gamma)
		c <- lab
	}

	// calculate minus and plus matrix
	minus(minusElm, elmIndex)
	plus(plusElm, elmIndex)
	delta := plusElm[elmIndex] - minusElm[elmIndex]

	// simulate image
	cCenterCalc := make(chan map[int][]float64)
	cMinusCalc := make(chan map[int][]float64)
	cPlusCalc := make(chan map[int][]float64)

	go simulateImageFromMat(cCenterCalc, elm)
	go simulateImageFromMat(cMinusCalc, minusElm)
	go simulateImageFromMat(cPlusCalc, plusElm)

	centerLabData := <-cCenterCalc
	minusLabData := <-cMinusCalc
	plusLabData := <-cPlusCalc

	// calculate all channel sensitivity
	var errorDivAve float64
	for index := 1; index < len(op.refLabData)+1; index++ {
		deltaERef2Center := calculator.DeltaECalculator(op.refLabData[index], centerLabData[index], op.kvals)
		deltaERef2Minus := calculator.DeltaECalculator(op.refLabData[index], minusLabData[index], op.kvals)
		deltaERef2Plus := calculator.DeltaECalculator(op.refLabData[index], plusLabData[index], op.kvals)

		div := (deltaERef2Plus - deltaERef2Minus) / delta
		errorDiv := (deltaERef2Center - target) * div

		errorDivAve += errorDiv
	}

	// calculate average
	errorDivAve /= float64(len(op.refCCData))

	// return
	return errorDivAve
}

/*
SetOptimizerCondition :Condition setter
*/
func (op *optimizer) SetOptimizerCondition(bachSize, maxTrial int, rangePercent float64) (bool, error) {

	// check input
	if bachSize < 1 {
		return false, errors.New("Bach size is <1, please correct bach size")
	}

	if maxTrial < 1 {
		return false, errors.New("Max trial number is <1, please correct bach size")
	}

	// update conditions
	op.condition.bachSize = bachSize
	op.condition.maxTrial = maxTrial
	op.condition.rangePercent = rangePercent / 100.0

	return true, nil
}

/*
SetOptimizerCoeffecients :Adam coeffecient setter
*/
func (op *optimizer) SetOptimizerCoeffecients(alpha, beta1, beta2, epsilon float64) (bool, error) {
	op.adamCondition.alpha = alpha
	op.adamCondition.beta1 = beta1
	op.adamCondition.beta2 = beta2
	op.adamCondition.epsilon = epsilon

	return true, nil
}

/*
OptimizedResults :Return the optimized result
*/
func (op *optimizer) OptimizedResults() (ccData map[int][]uint8, linearMat map[string]float64, redGain, blueGain float64, hueAngle int, brightness, saturation, deltaEAve float64) {
	return op.resultCCData, op.resultLinarMat, op.redGain, op.blueGain, op.resultHueAngle, op.resultBrightness, op.resultSaturation, op.resultDeltaE
}

/*
OptimizedEachPatchDeltaE :Return the optimized each pathc deltaE
*/
func (op *optimizer) OptimizedEachPatchDeltaE() map[int]float64 {
	return op.resultEachDeltaE
}

/*
Run :Run optimizer
*/
func (op *optimizer) Run(targetDeltaE float64) bool {
	elm := make([]float64, len(op.initLinearMatElm))
	copy(elm, op.initLinearMatElm)

	// randamly change elements
	randLinearMat := func() {
		// make rand seed from time
		rand.Seed(time.Now().UnixNano())

		// calculate value
		elmindex := rand.Intn(5)
		newValue := elm[elmindex] + rand.Float64()*op.condition.rangePercent*0.1
		if newValue > 0.001 || newValue < 1.0 {
			elm[elmindex] = newValue
		}
	}

	// deltaE calculation
	deltaEAve := func(data map[int][]uint8) float64 {
		labData := calculator.RGB2Lab(data, op.colorSpace, op.whitePoint, op.gamma)

		// calculate deltaE
		errAve := 0.
		for index := 1; index <= len(data); index++ {
			deltaE := calculator.DeltaECalculator(op.refLabData[index], labData[index], op.kvals)
			errAve += deltaE
		}
		// finalize calcualtion
		errAve /= float64(len(data))

		// return
		return errAve
	}

	/*
		Step-1 Minimize Delta-E by changing each elements value
	*/
	// minimize deltaE
	for epoc := 0; epoc < op.condition.maxTrial; epoc++ {
		// random select
		randLinearMat()

		for index := 0; index < len(elm); index++ {
			// initialize feedback factors
			m := 0.0
			v := 0.0

			// start batch
			for trial := 1; trial < op.condition.bachSize; trial++ {

				divDeltaE := op.gradient(elm, index, targetDeltaE)

				// calculate adam
				m = op.adamCondition.beta1*m + (1.0-op.adamCondition.beta1)*divDeltaE
				v = op.adamCondition.beta2*v + (1.0-op.adamCondition.beta2)*divDeltaE*divDeltaE

				mhat := m / (1.0 - math.Pow(op.adamCondition.beta1, float64(trial)))
				vhat := v / (1.0 - math.Pow(op.adamCondition.beta2, float64(trial)))

				val := -op.adamCondition.alpha * mhat / (math.Sqrt(vhat) + op.adamCondition.epsilon)
				updatedValue := elm[index] + val
				if updatedValue < 0.001 || updatedValue > 0.6 {
					updatedValue = elm[index]
				}

				// update
				elm[index] = updatedValue
			}
		}
	}

	// update Step-1 result
	op.resultCCData, op.redGain, op.blueGain = op.simulateImage(mapLinearMatElm(elm), op.gamma)
	op.resultDeltaE = deltaEAve(op.resultCCData)
	op.resultLinarMat = mapLinearMatElm(elm)

	// debug
	log.Println("--- Finished Linear Matrix Optimization ---")
	log.Println("Optimized Linear Matrix", op.resultLinarMat)
	log.Println("delta-E:", op.resultDeltaE)

	/*
		Step-2  Minimize Delta-E by roatating Hue
	*/
	op.resultHueAngle = 0
	for angle := op.condition.hue.minAngle; angle <= op.condition.hue.maxAngle; angle++ {

		buffer := op.simulation.AdjustHue(op.resultCCData, angle)
		errAve := deltaEAve(buffer)

		// update the best condition
		if errAve < op.resultDeltaE {
			op.resultDeltaE = errAve
			op.resultHueAngle = angle
			op.resultCCData = buffer
		}
	}

	log.Println("--- Finished Hue angle optimization ---")
	log.Println("Angle :", op.resultHueAngle, "Delta-E", op.resultDeltaE)

	/*
		Step-3 Brightness optimization
	*/
	for brightness := op.condition.brightness.low; brightness <= op.condition.brightness.high; brightness += 0.01 {

		buffer := op.simulation.AdjustBrightness(op.resultCCData, brightness)
		errAve := deltaEAve(buffer)

		// update the best condition
		if errAve < op.resultDeltaE {
			op.resultDeltaE = errAve
			op.resultBrightness = brightness
			op.resultCCData = buffer
		}
	}
	// debug
	log.Println("Optimized Brightness", op.resultBrightness)
	log.Println("delta-E:", op.resultDeltaE)

	/*
		Step-4 Saturation optimization
	*/
	for saturation := op.condition.brightness.low; saturation <= op.condition.brightness.high; saturation += 0.01 {

		buffer := op.simulation.AdjustSaturation(op.resultCCData, saturation)
		errAve := deltaEAve(buffer)

		// update the best condition
		if errAve < op.resultDeltaE {
			op.resultDeltaE = errAve
			op.resultSaturation = saturation
			op.resultCCData = buffer
		}
	}

	// debug
	log.Println("Optimized Saturation", op.resultSaturation)
	log.Println("delta-E:", op.resultDeltaE)

	// finalize all calcualtion
	op.resultEachDeltaE = make(map[int]float64, 0)
	labData := calculator.RGB2Lab(op.resultCCData, op.colorSpace, op.whitePoint, op.gamma)

	// calculate deltaE
	for index := 1; index <= len(op.resultCCData); index++ {
		deltaE := calculator.DeltaECalculator(op.refLabData[index], labData[index], op.kvals)
		op.resultEachDeltaE[index] = deltaE
	}

	return true
}
