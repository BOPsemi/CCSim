package calculator

import (
	"fmt"
	"testing"
)

func Test_Calculator(t *testing.T) {
	// linear mat elem moc
	linearMat := make(map[string]float64, 0)
	linearMat["a"] = 0
	linearMat["b"] = 0
	linearMat["c"] = 0
	linearMat["d"] = 0
	linearMat["e"] = 0
	linearMat["f"] = 0

	// r, gr, gb, b moc
	r := []float64{1, 0, 0}
	gr := []float64{0, 2, 0}
	gb := []float64{0, 2, 0}
	b := []float64{0, 0, 3}

	result := LinearMatixCalculator(r, gr, gb, b, linearMat)

	fmt.Println(result)
}
