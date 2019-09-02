package calculator

func illuminationIntegration(ill map[int]float64, conv map[int]float64) map[int]float64 {
	illInt := make(map[int]float64, 0)

	var sum float64
	for wavelength, data := range ill {
		sum += data * conv[wavelength]
	}

	for wavelength, illumination := range ill {
		illInt[wavelength] = 100 * illumination / sum
	}

	return illInt
}
