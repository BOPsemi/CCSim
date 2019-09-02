package controllers

import (
	"PixCSim/models"
	"errors"
)

/*
SETUPFOLDERNAME :define setup file directory name
*/
const (
	SETUPFOLDERNAME = "setup/"
)

/*
Setup :setup object do simulation env setup
*/
type Setup interface {
	DeviceWavelengthResponse() (map[int][]float64, error)
	LinearMatrixEelemtns() map[string]float64
}

// structure
type setup struct {
	rawDeviceWavelengthResponse [][]string
	rawLinearMatrixElements     [][]string

	deviceWavelengthResponse *models.DeviceResponse
	linearMatrixElements     *models.LinearMatrixElements
}

/*
NewSetup :initializer of setup object
*/
func NewSetup() Setup {
	obj := new(setup)

	// read setup files
	obj.readSetupFiles()

	// init objects
	obj.initObjects()

	return obj
}

/*
DeviceWavelengthResponse :return dev wavelenth res
*/
func (st *setup) DeviceWavelengthResponse() (map[int][]float64, error) {
	if st.deviceWavelengthResponse == nil {
		err := errors.New("Object is nil, please check object initialize")
		return nil, err
	}

	return st.deviceWavelengthResponse.SerializedData(), nil
}

func (st *setup) LinearMatrixEelemtns() map[string]float64 {
	return st.linearMatrixElements.Data
}

// read setup files
func (st *setup) readSetupFiles() bool {

	// loading file list
	loadingFileList := getFileListIn(SETUPFOLDERNAME, "*.csv")

	// initialize channel
	cDevRes := make(chan [][]string)
	cLinMat := make(chan [][]string)

	// read CSV file
	go readCSVFile(loadingFileList[0], cDevRes)
	go readCSVFile(loadingFileList[1], cLinMat)

	// update
	st.rawDeviceWavelengthResponse = <-cDevRes
	st.rawLinearMatrixElements = <-cLinMat

	return true
}

// initObjects
func (st *setup) initObjects() bool {
	cDevRes := make(chan bool)
	cLinMat := make(chan bool)

	devResInit := func(c chan bool) {
		st.deviceWavelengthResponse = models.NewDeviceResponse(st.rawDeviceWavelengthResponse)
		c <- true
	}

	linMatInit := func(c chan bool) {
		st.linearMatrixElements = models.NewLinearMatrixElements(st.rawLinearMatrixElements)
		c <- true
	}

	go devResInit(cDevRes)
	go linMatInit(cLinMat)

	devResState := <-cDevRes
	linMatState := <-cLinMat

	if !devResState {
		return false
	}

	if !linMatState {
		return false
	}

	return true
}
