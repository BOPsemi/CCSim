package models

/*
IlluminantConv :illuminant conversion factor
*/
type IlluminantConv struct {
	Data map[int]float64
}

/*
NewIlluminantConv :initialize object
*/
func NewIlluminantConv(data [][]string) *IlluminantConv {
	// check
	if data == nil {
		return nil
	}

	illConv := make(map[int]float64, 0)
	for _, row := range data {
		wavelength := StrToInt(row[0])
		conversion := StrToFloat64(row[1])

		illConv[wavelength] = conversion
	}

	// initialize
	obj := new(IlluminantConv)
	obj.Data = make(map[int]float64, 0)

	// update
	obj.Data = illConv

	return obj
}
