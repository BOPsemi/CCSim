package models

// ColorSpaceEnum :define color space enum
type ColorSpaceEnum int

/*
CIE :CIE color cpace
*/
const (
	CIE ColorSpaceEnum = iota + 1
	NTSC
	SRGB
)
