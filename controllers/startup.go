package controllers

import (
	"PixCSim/models"
)

/*
DATAFOLDERNAME :data folder name
*/
const (
	LOADFOLDERNAME = "load/"
)

/*
StartUp :start up controller, the class summarize program start up functions
*/
type StartUp interface {
	SetIllumination(illumination models.IlluminationEnum) bool
	SetScanWavelength(start, stop, step int) bool

	ColorChartReflection() (map[int][]float64, error)
	StandardMacbethColorCode() map[int]models.ColorPatch
	IlluminationSpectrum(illumination models.IlluminationEnum) map[int]float64
	IlluminantConvFactor() map[int]float64
}

// structure
type startup struct {
	// raw data stockers
	rawRGB8bit        [][]string
	rawReflectance    [][]string
	rawIllumination   [][]string
	rawIlluminantConv [][]string

	//
	ColorCode8bit   *models.MacbethColorCode
	PatchReflection *models.PatchReflection
	Illumination    *models.Illumination
	IlluminantConv  *models.IlluminantConv

	illuminationSetting models.IlluminationEnum

	wavelengthScan struct {
		start int // scan start
		stop  int // scan stop
		step  int // scan step
	}
}

/*
NewStartUp :initializer of start up object
*/
func NewStartUp() StartUp {
	// initialize object
	obj := new(startup)

	// initialize working directory
	obj.readStartupFiles()

	// initialize model objects
	if !obj.initObjects() {
		return nil
	}

	// initialize variables, default setting
	obj.illuminationSetting = models.D65
	obj.wavelengthScan.start = 400
	obj.wavelengthScan.stop = 700
	obj.wavelengthScan.step = 5

	return obj
}

/*
SetIllumination :setter of illlumination
*/
func (st *startup) SetIllumination(illumination models.IlluminationEnum) bool {
	st.illuminationSetting = illumination
	return true
}

/*
StandardMacbethColorCode :return standard Macbeth 24 color code
*/
func (st *startup) StandardMacbethColorCode() map[int]models.ColorPatch {
	return st.ColorCode8bit.Data
}

/*
IlluminantConvFactor :return illuminant conversion factor
*/
func (st *startup) IlluminantConvFactor() map[int]float64 {
	return st.IlluminantConv.Data
}

/*
IlluminationSpectrum :return the specified spectrum value
*/
func (st *startup) IlluminationSpectrum(illumination models.IlluminationEnum) map[int]float64 {
	spectrum := make(map[int]float64, 0)
	for wavelength, intensity := range st.Illumination.Data {
		switch illumination {
		case models.A:
			spectrum[wavelength] = intensity.A
		case models.B:
			spectrum[wavelength] = intensity.B
		case models.C:
			spectrum[wavelength] = intensity.C
		case models.D50:
			spectrum[wavelength] = intensity.D50
		case models.D55:
			spectrum[wavelength] = intensity.D55
		case models.D65:
			spectrum[wavelength] = intensity.D65
		case models.D75:
			spectrum[wavelength] = intensity.D75
		case models.ID50:
			spectrum[wavelength] = intensity.ID50
		case models.ID65:
			spectrum[wavelength] = intensity.ID65
		}
	}

	return spectrum
}

/*
Set wavelength scan
*/
func (st *startup) SetScanWavelength(start, stop, step int) bool {
	st.wavelengthScan.start = start
	st.wavelengthScan.stop = stop
	st.wavelengthScan.step = step

	return true
}

// read setup files
func (st *startup) readStartupFiles() bool {
	// list up the loding file
	loadingFileList := getFileListIn(LOADFOLDERNAME, "*.csv")

	// initialize channel
	cReflectance := make(chan [][]string)
	cRGB8bit := make(chan [][]string)
	cIllumination := make(chan [][]string)
	cIlluminantConv := make(chan [][]string)

	// read CSV file

	/*
		TODO :how to control order in the list
	*/
	go readCSVFile(loadingFileList[0], cRGB8bit)        // reflectance
	go readCSVFile(loadingFileList[1], cReflectance)    // RGB 8bit data
	go readCSVFile(loadingFileList[2], cIllumination)   // Illumination
	go readCSVFile(loadingFileList[3], cIlluminantConv) // Illumination Conv

	st.rawReflectance = <-cReflectance
	st.rawRGB8bit = <-cRGB8bit
	st.rawIllumination = <-cIllumination
	st.rawIlluminantConv = <-cIlluminantConv

	return true
}

// make object
func (st *startup) initObjects() bool {

	// define channel
	cCC8bit := make(chan bool)
	cIll := make(chan bool)
	cPatchRef := make(chan bool)
	cIllConv := make(chan bool)

	cc8BitInit := func(c chan bool) {
		st.ColorCode8bit = models.NewMacbethColorCode(st.rawRGB8bit)
		c <- true
	}

	illInit := func(c chan bool) {
		st.Illumination = models.NewIllumination(st.rawIllumination)
		c <- true
	}

	patchRefInit := func(c chan bool) {
		st.PatchReflection = models.NewPatchReflection(st.rawReflectance)
		c <- true
	}

	illConvInit := func(c chan bool) {
		st.IlluminantConv = models.NewIlluminantConv(st.rawIlluminantConv)
		c <- true
	}

	go cc8BitInit(cCC8bit)
	go illInit(cIll)
	go patchRefInit(cPatchRef)
	go illConvInit(cIllConv)

	cc8BitState := <-cCC8bit
	illState := <-cIll
	patchRefState := <-cPatchRef
	illConvState := <-cIllConv

	if !(cc8BitState && illState && patchRefState && illConvState) {
		return false
	}

	return true
}

/*
PatchReflectionAt :return at the calcualted reflection under the specified illumination
*/
func (st *startup) patchReflectionAt(wavelength int, illumination models.IlluminationEnum) []float64 {

	// inline macro to calculat intensity
	calculateEachPatchReflectance := func(illumination float64) []float64 {
		patchReflections := make([]float64, 0)
		for _, patch := range st.PatchReflection.Data {
			eachPatchReflection := patch[wavelength] * illumination
			patchReflections = append(patchReflections, eachPatchReflection)
		}
		return patchReflections
	}

	// calculate each patch reflectance
	switch illumination {
	case models.A:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].A)
	case models.B:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].B)
	case models.C:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].C)
	case models.D50:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].D50)
	case models.D55:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].D55)
	case models.D65:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].D65)
	case models.D75:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].D75)
	case models.ID50:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].ID50)
	case models.ID65:
		return calculateEachPatchReflectance(st.Illumination.Data[wavelength].ID65)
	default:
		return nil
	}
}

/*
ColorChartReflection :calculte each color chart reflection and scan wavelength
*/
func (st *startup) ColorChartReflection() (map[int][]float64, error) {
	stocker := make(map[int][]float64, 0)

	for wavelength := st.wavelengthScan.start; wavelength <= st.wavelengthScan.stop; wavelength += st.wavelengthScan.step {
		stocker[wavelength] = st.patchReflectionAt(wavelength, st.illuminationSetting)
	}

	return stocker, nil
}
