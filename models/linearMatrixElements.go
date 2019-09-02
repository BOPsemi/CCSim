package models

/*
LinearMatrixElements :linear matrix elements
*/
type LinearMatrixElements struct {
	Data map[string]float64
}

/*
NewLinearMatrixElements :initialize linear matrix
*/
func NewLinearMatrixElements(data [][]string) *LinearMatrixElements {
	if data == nil {
		return nil
	}

	obj := new(LinearMatrixElements)
	obj.Data = make(map[string]float64, 0)

	var elmStocker []float64
	for _, row := range data {
		elm := StrToFloat64(row[0])
		elmStocker = append(elmStocker, elm)
	}

	// update data
	obj.Data["a"] = elmStocker[0]
	obj.Data["b"] = elmStocker[1]
	obj.Data["c"] = elmStocker[2]
	obj.Data["d"] = elmStocker[3]
	obj.Data["e"] = elmStocker[4]
	obj.Data["f"] = elmStocker[5]

	// return
	return obj
}
