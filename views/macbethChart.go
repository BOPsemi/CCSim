package views

import (
	"image"
	"os"

	"github.com/fogleman/gg"
)

/*
OUTPUTFOLDER :output folder name
*/
const (
	OUTPUTFOLDER = "output/"
	LEFTMARGIN   = 30
	RIGHTMARGIN  = 30
	TOPMARGIN    = 30
	BOTTOMMARGIN = 30
	SPACING      = 10

	PATCHWIDTH  = 100
	PATCHHEIGHT = 100
)

/*
MacbethChart :interface of Macbeth chart
*/
type MacbethChart interface {
	// margin and patch size setter
	SetPatchMargin(margins ...map[string]int) bool
	SetPatchSize(size ...map[string]int) bool

	// create image
	CreateMachbethChart(data map[int][]uint8) bool
	Save(pngFileName ...string) bool
}

type macbethChart struct {
	outputDir    string // output directory
	patchImages  map[int]image.Image
	macbethChart image.Image

	// patch size
	patchSize struct {
		height int
		width  int
	}

	// patch margin
	margin struct {
		top     int
		bottom  int
		right   int
		left    int
		spacing int
	}
}

/*
NewMacbethChart :initializer of MacbethChart
*/
func NewMacbethChart(outputdir ...string) MacbethChart {
	obj := new(macbethChart)

	// get working dir
	workDir := os.Getenv("APPATH")

	// output directory
	if outputdir == nil {
		obj.outputDir = workDir + OUTPUTFOLDER
	} else {
		obj.outputDir = workDir + outputdir[0]
	}

	// initialize magines
	obj.margin.top = TOPMARGIN
	obj.margin.bottom = BOTTOMMARGIN
	obj.margin.right = RIGHTMARGIN
	obj.margin.left = LEFTMARGIN
	obj.margin.spacing = SPACING

	// initialize patch size
	obj.patchSize.height = PATCHHEIGHT
	obj.patchSize.width = PATCHWIDTH

	return obj
}

// setter for patch margin
func (mc *macbethChart) SetPatchMargin(margins ...map[string]int) bool {
	if margins != nil {
		for _, setMargin := range margins {
			mc.margin.top = setMargin["top"]
			mc.margin.bottom = setMargin["bottom"]
			mc.margin.left = setMargin["left"]
			mc.margin.right = setMargin["right"]
			mc.margin.spacing = setMargin["spacing"]
		}
	}

	return true
}

// setter for patch size
func (mc *macbethChart) SetPatchSize(size ...map[string]int) bool {
	if size != nil {
		for _, patchSize := range size {
			mc.patchSize.height = patchSize["height"]
			mc.patchSize.width = patchSize["width"]
		}
	}
	return true
}

/*
CreateMachbethChart :function of 24 color patch creation
*/
func (mc *macbethChart) CreateMachbethChart(data map[int][]uint8) bool {

	// inline macro
	// extract rgb data from slice
	extractRGB := func(rgbData []uint8) (r, g, b uint8) {
		return rgbData[0], rgbData[1], rgbData[2]
	}

	// patch generate
	patchImageGenerator := func(r, g, b uint8) image.Image {
		context := gg.NewContext(mc.patchSize.width, mc.patchSize.height)
		context.DrawRectangle(0.0, 0.0, float64(mc.patchSize.width), float64(mc.patchSize.height))
		context.SetRGBA255(int(r), int(g), int(b), 255)
		context.Fill()

		return context.Image()
	}

	// initialize patch images
	mc.patchImages = make(map[int]image.Image, 0)
	for index, rgbData := range data {
		mc.patchImages[index] = patchImageGenerator(extractRGB(rgbData))
	}

	// create 24 patch image on context
	// extract one patch size

	bkwidth := mc.margin.left + 6*mc.patchSize.width + 5*mc.margin.spacing + mc.margin.right
	bkheight := mc.margin.top + 4*mc.patchSize.height + 3*mc.margin.spacing + mc.margin.bottom

	// make background context
	context := gg.NewContext(bkwidth, bkheight)
	context.SetRGBA255(0, 0, 0, 255)
	context.DrawRectangle(0.0, 0.0, float64(bkwidth), float64(bkheight))
	context.Fill()

	// Put each patch image
	for index := 1; index < len(mc.patchImages)+1; index++ {
		data := mc.patchImages[index]
		// point data
		if index <= 6 {
			// 1st row
			i := index - 1
			xpoint := mc.margin.left + i*mc.patchSize.width + i*mc.margin.spacing
			ypoint := mc.margin.top

			context.DrawImage(data, xpoint, ypoint)

		} else if index > 5 && index <= 12 {
			// 2nd row
			i := (index - 7)
			xpoint := mc.margin.left + i*mc.patchSize.width + i*mc.margin.spacing
			ypoint := mc.margin.top + mc.patchSize.height + mc.margin.spacing

			context.DrawImage(data, xpoint, ypoint)

		} else if index > 10 && index <= 18 {
			// 3rd row
			i := (index - 13)
			xpoint := mc.margin.left + i*mc.patchSize.width + i*mc.margin.spacing
			ypoint := mc.margin.top + 2*mc.patchSize.height + 2*mc.margin.spacing

			context.DrawImage(data, xpoint, ypoint)

		} else {
			// 4th row
			i := (index - 19)
			xpoint := mc.margin.left + i*mc.patchSize.width + i*mc.margin.spacing
			ypoint := mc.margin.top + 3*mc.patchSize.height + 3*mc.margin.spacing

			context.DrawImage(data, xpoint, ypoint)
		}
	}
	// update
	mc.macbethChart = context.Image()

	// return
	return true
}

/*
Save png file
*/
func (mc *macbethChart) Save(pngFileName ...string) bool {
	// file name check
	var fileName string
	if pngFileName == nil {
		fileName = "macbeth_24Chart"
	} else {
		fileName = pngFileName[0]
	}

	// save image
	status, err := saveImageAsPNG(mc.outputDir, fileName, mc.macbethChart)
	if err != nil {
		return false
	}
	return status
}
