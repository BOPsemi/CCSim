package calculator

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// for debugging
func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}

/*
LinearMatixCalculator :calculate linear matrix
*/
func LinearMatixCalculator(r, gr, gb, b []float64, linearMatElm map[string]float64) map[int][]float64 {
	// stocker of result
	resultMap := make(map[int][]float64, 0)

	// normarize elements
	rgain := 1.0 + linearMatElm["a"] + linearMatElm["b"]
	ggain := 1.0 + linearMatElm["c"] + linearMatElm["d"]
	bgain := 1.0 + linearMatElm["e"] + linearMatElm["f"]

	// make linear matrix from slice
	linearMatElmSlice := []float64{
		1.0, -linearMatElm["a"] / rgain, -linearMatElm["b"] / rgain,
		-linearMatElm["c"] / ggain, 1.0, -linearMatElm["d"] / ggain,
		-linearMatElm["e"] / bgain, -linearMatElm["f"] / bgain, 1.0,
	}
	linearMat := mat.NewDense(3, 3, linearMatElmSlice)

	// calcualte each patch
	for index := 0; index < len(r); index++ {
		// calculation result
		resultMat := mat.NewDense(3, 1, nil)

		// make rgb matrix
		rsignal := r[index]
		gsignal := (gr[index] + gb[index]) / 2.0
		bsignal := b[index]

		// rgb matrix
		rgb := mat.NewDense(3, 1, []float64{rsignal, gsignal, bsignal})

		// calculation
		resultMat.Mul(linearMat, rgb)

		// result slice
		resultMatSlice := []float64{resultMat.At(0, 0), resultMat.At(1, 0), resultMat.At(2, 0)}

		// update
		resultMap[index+1] = resultMatSlice
	}

	return resultMap
}
