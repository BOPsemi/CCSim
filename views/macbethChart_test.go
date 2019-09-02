package views

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mocData() map[int][]uint8 {
	data := make(map[int][]uint8, 0)

	data[1] = []uint8{100, 150, 200}

	return data
}

func mocArgs() map[string]int {
	data := make(map[string]int, 0)

	data["top"] = 122
	data["left"] = 222

	return data
}

func TestNewMacbethChart(t *testing.T) {
	//inputData := mocData()

	obj := NewMacbethChart()
	assert.NotNil(t, obj)
}

func TestSetPatchMargin(t *testing.T) {
	//inputData := mocData()
	inputArgs := mocArgs()

	obj := NewMacbethChart()
	obj.SetPatchMargin(inputArgs)
}

func TestComFunc(t *testing.T) {
	str := checkLastWordInPath("hoge/")
	fmt.Println(str)
}
