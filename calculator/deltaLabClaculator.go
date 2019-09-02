package calculator

import "math"

/*
DeltaECalculator :calculate delta Lab, calculation algorithm is CIE2000
*/
func DeltaECalculator(ref []float64, comp []float64, kvals []float64) float64 {
	// check input
	if ref == nil || comp == nil {
		return -999
	}

	// inline
	rad2Deg := func(rad float64) float64 {
		return 180 * rad / math.Pi
	}

	deg2Rad := func(deg float64) float64 {
		return deg * math.Pi / 180.0
	}

	/*
		[0]	:L
		[1]	:a
		[2]	:b
	*/
	deltaLp := comp[0] - ref[0]
	lAve := (comp[0] + ref[0]) / 2.0

	// calculate cAve
	c1 := math.Sqrt(math.Pow(ref[1], 2.0) + math.Pow(ref[2], 2.0))
	c2 := math.Sqrt(math.Pow(comp[1], 2.0) + math.Pow(comp[2], 2.0))
	cAve := (c1 + c2) / 2.0

	ap1 := ref[1] + (ref[1]/2.0)*(1.0-math.Sqrt(math.Pow(cAve, 7.0)/(math.Pow(cAve, 7.0)+math.Pow(25, 7))))
	ap2 := comp[1] + (comp[1]/2.0)*(1.0-math.Sqrt(math.Pow(cAve, 7.0)/(math.Pow(cAve, 7.0)+math.Pow(25, 7))))

	cp1 := math.Sqrt(math.Pow(ap1, 2.0) + math.Pow(ref[2], 2.0))
	cp2 := math.Sqrt(math.Pow(ap2, 2.0) + math.Pow(comp[2], 2.0))

	cpAve := (cp1 + cp2) / 2.0
	deltaCp := cp2 - cp1

	var hp1 float64
	if ref[2] == 0 && ap1 == 0 {
		hp1 = 0.0
	} else {
		hp1 = rad2Deg(math.Atan2(ref[2], ap1))
		if hp1 < 0 {
			hp1 += 360.0
		}
	}

	var hp2 float64
	if ref[2] == 0 && ap1 == 0 {
		hp2 = 0.0
	} else {
		hp2 = rad2Deg(math.Atan2(comp[2], ap2))
		if hp2 < 0 {
			hp2 += 360.0
		}
	}

	var deltahp float64
	if c1 == 0.0 || c2 == 0.0 {
		deltahp = 0.0
	} else if math.Abs(hp1-hp2) <= 180.0 {
		deltahp = hp2 - hp1
	} else if hp2 <= hp1 {
		deltahp = hp2 - hp1
	} else {
		deltahp = hp2 - hp1 - 360.0
	}
	deltaHp := 2.0 * math.Sqrt(cp1*cp2) * math.Sin(deg2Rad(deltahp)/2.0)

	var HpAve float64
	if math.Abs(hp1-hp2) > 180.0 {
		HpAve = (hp1 + hp2 + 360.0) / 2.0
	} else {
		HpAve = (hp1 + hp2 + 360.0) / 2.0
	}

	t := 1.0 -
		0.17*math.Cos(deg2Rad(HpAve-30.0)) +
		0.24*math.Cos(deg2Rad(2.0*HpAve)) +
		0.32*math.Cos(deg2Rad(3.0*HpAve+6.0)) -
		0.20*math.Cos(deg2Rad(4.0*HpAve-63.0))

	sl := 1.0 + ((0.015 * math.Pow(lAve-50.0, 2.0)) / math.Sqrt(20.0+math.Pow(lAve-50.0, 2.0)))
	sc := 1.0 + 0.045*cpAve
	sh := 1.0 + 0.015*cpAve*t

	rt := -2.0 * math.Sqrt(math.Pow(cpAve, 7.0)/(math.Pow(cpAve, 7.0)+math.Pow(25.0, 7.0))) *
		math.Sin(deg2Rad(60.0*math.Exp(-math.Pow((HpAve-275.0)/25.0, 2.0))))

	result := math.Sqrt(
		math.Pow(deltaLp/(kvals[0]*sl), 2.0) +
			math.Pow(deltaCp/(kvals[1]*sc), 2.0) +
			math.Pow(deltaHp/(kvals[2]*sh), 2.0) +
			rt*(deltaCp/(kvals[1]*sc))*(deltaHp/(kvals[2]*sh)))

	return result
}
