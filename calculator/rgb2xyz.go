package calculator

import (
	"PixCSim/models"

	"gonum.org/v1/gonum/mat"
)

/*
ConversionMatrixInit :calculate conversion matrix from color space and gamma
*/
func ConversionMatrixInit(cspace models.ColorSpaceEnum, wpoint *models.WhitePoint) []float64 {
	/*
		Referece link
		Basically, just copied this calculation in this part
		http://technorgb.blogspot.com/2015/08/rgb-xyz.html
	*/

	// initialize white point
	var rgbElm, wElm []float64
	//wpoint := models.NewWhitePoint(cspace)

	switch cspace {
	case models.CIE:
		rgbElm = wpoint.CIE.RGBElm
		wElm = wpoint.CIE.WhiteElm

	case models.NTSC:
		rgbElm = wpoint.NTSC.RGBElm
		wElm = wpoint.NTSC.WhiteElm

	default:
		rgbElm = wpoint.SRGB.RGBElm
		wElm = wpoint.SRGB.WhiteElm
	}

	// create matrix from white point data
	rgbElmMat := mat.NewDense(3, 3, rgbElm)
	wElmMat := mat.NewDense(3, 1, wElm)

	// initialize buffer
	invRGBElmMat := mat.NewDense(3, 3, nil)
	wmat2 := mat.NewDense(3, 1, nil)
	resultMat := mat.NewDense(3, 3, nil)

	// calculate inverse matrix
	invRGBElmMat.Inverse(rgbElmMat)

	// M^-1 x w
	wmat2.Mul(invRGBElmMat, wElmMat)

	// create transversal rgb matrix
	/*
		|	fx	0	0 	|
		|	0	fy	0	|
		|	0	0	fz	|
	*/
	trgb := []float64{
		wmat2.At(0, 0), 0.0, 0.0,
		0.0, wmat2.At(1, 0), 0.0,
		0.0, 0.0, wmat2.At(2, 0),
	}
	trgbMat := mat.NewDense(3, 3, trgb)

	// result
	resultMat.Mul(rgbElmMat, trgbMat) // rgbElm x trgb

	resultMatElm := []float64{
		resultMat.At(0, 0), resultMat.At(0, 1), resultMat.At(0, 2),
		resultMat.At(1, 0), resultMat.At(1, 1), resultMat.At(1, 2),
		resultMat.At(2, 0), resultMat.At(2, 1), resultMat.At(2, 2),
	}

	return resultMatElm
}

/*
RGB2XYZ :RGB to XYZ converter
*/
func RGB2XYZ(eachPatchRGBdata map[int][]uint8, gamma float64, convMatElm []float64) map[int][]float64 {

	// calculate each patch XYZ
	eachPatchDegammaRGBdata := DeGammaCorrection(eachPatchRGBdata, gamma)

	// initialize result
	xyzResults := make(map[int][]float64, 0)

	// convert each patch data from rgb to xyz
	for index, rgbData := range eachPatchDegammaRGBdata {
		/*
			Calculate XYZ signal
		*/
		rgbSignalMat := mat.NewDense(3, 1, rgbData)
		elmMatrix := mat.NewDense(3, 3, convMatElm)
		resultMatrix := mat.NewDense(3, 1, nil)

		resultMatrix.Mul(elmMatrix, rgbSignalMat)
		xyzResult := []float64{
			resultMatrix.At(0, 0),
			resultMatrix.At(1, 0),
			resultMatrix.At(2, 0),
		}

		// x, y, z
		xyzResults[index] = xyzResult
	}

	return xyzResults
}
