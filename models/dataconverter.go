package models

import "strconv"

/*
StrToInt :Define String to Int converter
*/
func StrToInt(str string) int {
	// define string to Int inline function
	value, err := strconv.Atoi(str)
	if err != nil {
		return -999
	}
	return value
}

/*
StrToUint8 :Define String to UInt8 converter
*/
func StrToUint8(str string) uint8 {
	value := StrToInt(str)
	if value == -999 {
		return 255
	}

	if value < 0 {
		value = 0
	}

	if value > 255 {
		value = 255
	}

	return uint8(value)
}

/*
StrToFloat64 :Define String to Float64 converter
*/
func StrToFloat64(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -999.9
	}

	return value
}
