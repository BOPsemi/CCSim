package controllers

import (
	"PixCSim/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUp() StartUp {
	return NewStartUp()
}

func Test_StartUp(t *testing.T) {
	obj := setUp()
	assert.NotNil(t, obj)

	data, err := obj.ColorChartReflection()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(data))
}

func Test_SetUp(t *testing.T) {
	obj := NewSetup()
	assert.NotNil(t, obj)

	data, _ := obj.DeviceWavelengthResponse()
	fmt.Println(data)

}

func Test_IlluminantConvFactor(t *testing.T) {
	obj := NewStartUp()
	assert.NotNil(t, obj)

	/*
		illConv := obj.IlluminantConvFactor()
		fmt.Println(illConv)
	*/

	illSpectrum := obj.IlluminationSpectrum(models.D65)
	fmt.Println(illSpectrum)

}
