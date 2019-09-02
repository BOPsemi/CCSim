package models

/*
IlluminationIntensity :intensity of illumination
*/
type IlluminationIntensity struct {
	A    float64
	B    float64
	C    float64
	D50  float64
	D55  float64
	D65  float64
	D75  float64
	ID50 float64
	ID65 float64
}

/*
Illumination :definition of illumination
*/
type Illumination struct {
	Data map[int]IlluminationIntensity // wavelength and intensity
}

/*
NewIllumination :initializer of Illumination object
*/
func NewIllumination(data [][]string) *Illumination {
	if data == nil {
		return nil
	}

	// initialize object
	obj := new(Illumination)
	obj.Data = make(map[int]IlluminationIntensity, 0)

	for _, row := range data {
		// pick up row data
		wavelength := StrToInt(row[0])
		illObj := &IlluminationIntensity{
			A:    StrToFloat64(row[1]),
			B:    StrToFloat64(row[2]),
			C:    StrToFloat64(row[3]),
			D50:  StrToFloat64(row[4]),
			D55:  StrToFloat64(row[5]),
			D65:  StrToFloat64(row[6]),
			D75:  StrToFloat64(row[7]),
			ID50: StrToFloat64(row[8]),
			ID65: StrToFloat64(row[9]),
		}

		// update object
		obj.Data[wavelength] = *illObj
	}
	return obj
}
