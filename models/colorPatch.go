package models

/*
WHITEDIGITALNUMBER : patch 19 digital number
*/
const (
	WHITEDIGITALNUMBER = 243
)

/*
RGBChannel :defintion of color channel
*/
type RGBChannel struct {
	R uint8 // Red channel
	G uint8 // Green channel
	B uint8 // Blue channel
}

/*
MacbethColorCode :24 color Macbeth Color Patch
*/
type MacbethColorCode struct {
	Data map[int]ColorPatch
}

/*
NewMacbethColorCode :initializer of MacbethColorPatch object
*/
func NewMacbethColorCode(data [][]string) *MacbethColorCode {
	if data == nil {
		return nil
	}

	obj := new(MacbethColorCode)
	obj.Data = make(map[int]ColorPatch, 0)

	for index, patch := range data {
		ccPatch := NewColorPatch(patch)

		obj.Data[index+1] = *ccPatch
	}

	return obj
}

/*
NewMacbethColorCodeFromRaw :make object from raw data
*/
func NewMacbethColorCodeFromRaw(data map[int][]float64) *MacbethColorCode {
	if data == nil {
		return nil
	}

	// create object
	obj := new(MacbethColorCode)
	obj.Data = make(map[int]ColorPatch, 0)

	// define degitizer
	digitizer := func(value float64) uint8 {
		// binning < 0
		if value < 0 {
			return 0
		}

		// binning > 255
		digitalNumber := value * WHITEDIGITALNUMBER
		if digitalNumber > 255 {
			return 255
		}

		return uint8(digitalNumber)
	}

	// mapping
	for index, rgbData := range data {

		channel := &RGBChannel{
			R: digitizer(rgbData[0]),
			G: digitizer(rgbData[1]),
			B: digitizer(rgbData[2]),
		}

		ccPatch := &ColorPatch{
			PatchNumber: index,
			PatchName:   "",
			Channel:     channel,
		}

		obj.Data[index] = *ccPatch
	}

	return obj

}

/*
ColorPatch :definition of color patch data
*/
type ColorPatch struct {
	PatchNumber int
	PatchName   string
	Channel     *RGBChannel
}

/*
NewColorPatch :initializer of ColorPatch
*/
func NewColorPatch(data []string) *ColorPatch {
	// initialize channel
	channel := &RGBChannel{
		R: StrToUint8(data[2]),
		G: StrToUint8(data[3]),
		B: StrToUint8(data[4]),
	}

	// return object
	return &ColorPatch{
		PatchNumber: StrToInt(data[0]),
		PatchName:   data[1],
		Channel:     channel,
	}
}
