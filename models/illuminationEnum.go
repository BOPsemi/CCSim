package models

/*
IlluminationEnum :definition of illumination Enum
*/
type IlluminationEnum int

/*
A :Illumination-A
*/
const (
	A IlluminationEnum = iota + 1
	B
	C
	D50
	D55
	D65
	D75
	ID50
	ID65
)
