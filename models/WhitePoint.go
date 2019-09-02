package models

// RGBWP :RGB white point
type RGBWP struct {
	X float64
	Y float64
	Z float64
}

// WhiteWP :White white point
type WhiteWP struct {
	X2 float64
	Y2 float64
	Z2 float64
}

// Points :Point
type Points struct {
	R  *RGBWP
	G  *RGBWP
	B  *RGBWP
	W  *WhiteWP
	Wn struct {
		X float64
		Y float64
		Z float64
	}
	Wns      []float64
	RGBElm   []float64
	WhiteElm []float64
}

// WhitePoint :final one
type WhitePoint struct {
	CIE  *Points
	NTSC *Points
	SRGB *Points
}

/*
NewWhitePoint :initialize white point object
*/
func NewWhitePoint(cspace ColorSpaceEnum) *WhitePoint {

	obj := new(WhitePoint)

	// CIE
	obj.CIE = &Points{
		R: &RGBWP{
			X: 0.735,
			Y: 0.265,
		},

		G: &RGBWP{
			X: 0.274,
			Y: 0.717,
		},

		B: &RGBWP{
			X: 0.167,
			Y: 0.009,
		},

		W: &WhiteWP{
			X2: 0.3333,
			Y2: 0.3333,
		},
	}

	// NTSC
	obj.NTSC = &Points{
		R: &RGBWP{
			X: 0.67,
			Y: 0.33,
		},

		G: &RGBWP{
			X: 0.21,
			Y: 0.71,
		},

		B: &RGBWP{
			X: 0.14,
			Y: 0.08,
		},

		W: &WhiteWP{
			X2: 0.31006,
			Y2: 0.31616,
		},
	}

	// sRGB
	obj.SRGB = &Points{
		R: &RGBWP{
			X: 0.64,
			Y: 0.33,
		},

		G: &RGBWP{
			X: 0.3,
			Y: 0.6,
		},

		B: &RGBWP{
			X: 0.15,
			Y: 0.06,
		},

		W: &WhiteWP{
			X2: 0.3127,
			Y2: 0.3290,
		},
	}

	obj.initWhitePointMatrix() // initialize White Point Matrix

	return obj
}

// initialize white pixel matrix
func (wp *WhitePoint) initWhitePointMatrix() bool {
	// define inline function
	rgbZ := func(x, y float64) float64 {
		return (1.0 - x - y)
	}
	rgbWn := func(x, y float64) float64 {
		return (x / y)
	}

	// CIE
	wp.CIE.R.Z = rgbZ(wp.CIE.R.X, wp.CIE.R.Y)
	wp.CIE.G.Z = rgbZ(wp.CIE.G.X, wp.CIE.G.Y)
	wp.CIE.B.Z = rgbZ(wp.CIE.B.X, wp.CIE.B.Y)
	wp.CIE.W.Z2 = rgbZ(wp.CIE.W.X2, wp.CIE.W.Y2)

	wp.CIE.Wn.X = rgbWn(wp.CIE.W.X2, wp.CIE.W.Y2)
	wp.CIE.Wn.Y = rgbWn(wp.CIE.W.Y2, wp.CIE.W.Y2)
	wp.CIE.Wn.Z = rgbWn(wp.CIE.W.Z2, wp.CIE.W.Y2)

	wp.CIE.RGBElm = []float64{
		wp.CIE.R.X, wp.CIE.G.X, wp.CIE.B.X,
		wp.CIE.R.Y, wp.CIE.G.Y, wp.CIE.B.Y,
		wp.CIE.R.Z, wp.CIE.G.Z, wp.CIE.B.Z,
	}

	wp.CIE.WhiteElm = []float64{
		wp.CIE.Wn.X, wp.CIE.Wn.Y, wp.CIE.Wn.Z,
	}

	// NTSC
	wp.NTSC.R.Z = rgbZ(wp.NTSC.R.X, wp.NTSC.R.Y)
	wp.NTSC.G.Z = rgbZ(wp.NTSC.G.X, wp.NTSC.G.Y)
	wp.NTSC.B.Z = rgbZ(wp.NTSC.B.X, wp.NTSC.B.Y)
	wp.NTSC.W.Z2 = rgbZ(wp.NTSC.W.X2, wp.NTSC.W.Y2)

	wp.NTSC.Wn.X = rgbWn(wp.NTSC.W.X2, wp.NTSC.W.Y2)
	wp.NTSC.Wn.Y = rgbWn(wp.NTSC.W.Y2, wp.NTSC.W.Y2)
	wp.NTSC.Wn.Z = rgbWn(wp.NTSC.W.Z2, wp.NTSC.W.Y2)

	wp.NTSC.RGBElm = []float64{
		wp.NTSC.R.X, wp.NTSC.G.X, wp.NTSC.B.X,
		wp.NTSC.R.Y, wp.NTSC.G.Y, wp.NTSC.B.Y,
		wp.NTSC.R.Z, wp.NTSC.G.Z, wp.NTSC.B.Z,
	}

	wp.NTSC.WhiteElm = []float64{
		wp.NTSC.Wn.X, wp.NTSC.Wn.Y, wp.NTSC.Wn.Z,
	}

	// sRGB
	wp.SRGB.R.Z = rgbZ(wp.SRGB.R.X, wp.SRGB.R.Y)
	wp.SRGB.G.Z = rgbZ(wp.SRGB.G.X, wp.SRGB.G.Y)
	wp.SRGB.B.Z = rgbZ(wp.SRGB.B.X, wp.SRGB.B.Y)
	wp.SRGB.W.Z2 = rgbZ(wp.SRGB.W.X2, wp.SRGB.W.Y2)

	wp.SRGB.Wn.X = rgbWn(wp.SRGB.W.X2, wp.SRGB.W.Y2)
	wp.SRGB.Wn.Y = rgbWn(wp.SRGB.W.Y2, wp.SRGB.W.Y2)
	wp.SRGB.Wn.Z = rgbWn(wp.SRGB.W.Z2, wp.SRGB.W.Y2)

	wp.SRGB.RGBElm = []float64{
		wp.SRGB.R.X, wp.SRGB.G.X, wp.SRGB.B.X,
		wp.SRGB.R.Y, wp.SRGB.G.Y, wp.SRGB.B.Y,
		wp.SRGB.R.Z, wp.SRGB.G.Z, wp.SRGB.B.Z,
	}

	wp.SRGB.WhiteElm = []float64{
		wp.SRGB.Wn.X, wp.SRGB.Wn.Y, wp.SRGB.Wn.Z,
	}

	return true
}
