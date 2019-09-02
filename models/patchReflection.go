package models

/*
PatchReflection :definition of color chart each patch reflection
*/
type PatchReflection struct {
	Data []map[int]float64
}

/*
NewPatchReflection :initializer of patch reflection
*/
func NewPatchReflection(data [][]string) *PatchReflection {
	if data == nil {
		return nil
	}

	obj := new(PatchReflection)
	obj.Data = make([]map[int]float64, 0)

	patch01 := make(map[int]float64, 0)
	patch02 := make(map[int]float64, 0)
	patch03 := make(map[int]float64, 0)
	patch04 := make(map[int]float64, 0)
	patch05 := make(map[int]float64, 0)
	patch06 := make(map[int]float64, 0)
	patch07 := make(map[int]float64, 0)
	patch08 := make(map[int]float64, 0)
	patch09 := make(map[int]float64, 0)
	patch10 := make(map[int]float64, 0)
	patch11 := make(map[int]float64, 0)
	patch12 := make(map[int]float64, 0)
	patch13 := make(map[int]float64, 0)
	patch14 := make(map[int]float64, 0)
	patch15 := make(map[int]float64, 0)
	patch16 := make(map[int]float64, 0)
	patch17 := make(map[int]float64, 0)
	patch18 := make(map[int]float64, 0)
	patch19 := make(map[int]float64, 0)
	patch20 := make(map[int]float64, 0)
	patch21 := make(map[int]float64, 0)
	patch22 := make(map[int]float64, 0)
	patch23 := make(map[int]float64, 0)
	patch24 := make(map[int]float64, 0)

	for _, row := range data {
		wavelenght := StrToInt(row[0])

		patch01[wavelenght] = StrToFloat64(row[1])
		patch02[wavelenght] = StrToFloat64(row[2])
		patch03[wavelenght] = StrToFloat64(row[3])
		patch04[wavelenght] = StrToFloat64(row[4])
		patch05[wavelenght] = StrToFloat64(row[5])
		patch06[wavelenght] = StrToFloat64(row[6])
		patch07[wavelenght] = StrToFloat64(row[7])
		patch08[wavelenght] = StrToFloat64(row[8])
		patch09[wavelenght] = StrToFloat64(row[9])
		patch10[wavelenght] = StrToFloat64(row[10])
		patch11[wavelenght] = StrToFloat64(row[11])
		patch12[wavelenght] = StrToFloat64(row[12])
		patch13[wavelenght] = StrToFloat64(row[13])
		patch14[wavelenght] = StrToFloat64(row[14])
		patch15[wavelenght] = StrToFloat64(row[15])
		patch16[wavelenght] = StrToFloat64(row[16])
		patch17[wavelenght] = StrToFloat64(row[17])
		patch18[wavelenght] = StrToFloat64(row[18])
		patch19[wavelenght] = StrToFloat64(row[19])
		patch20[wavelenght] = StrToFloat64(row[20])
		patch21[wavelenght] = StrToFloat64(row[21])
		patch22[wavelenght] = StrToFloat64(row[22])
		patch23[wavelenght] = StrToFloat64(row[23])
		patch24[wavelenght] = StrToFloat64(row[24])

	}

	// stacking
	obj.Data = append(obj.Data, patch01)
	obj.Data = append(obj.Data, patch02)
	obj.Data = append(obj.Data, patch03)
	obj.Data = append(obj.Data, patch04)
	obj.Data = append(obj.Data, patch05)
	obj.Data = append(obj.Data, patch06)
	obj.Data = append(obj.Data, patch07)
	obj.Data = append(obj.Data, patch08)
	obj.Data = append(obj.Data, patch09)
	obj.Data = append(obj.Data, patch10)
	obj.Data = append(obj.Data, patch11)
	obj.Data = append(obj.Data, patch12)
	obj.Data = append(obj.Data, patch13)
	obj.Data = append(obj.Data, patch14)
	obj.Data = append(obj.Data, patch15)
	obj.Data = append(obj.Data, patch16)
	obj.Data = append(obj.Data, patch17)
	obj.Data = append(obj.Data, patch18)
	obj.Data = append(obj.Data, patch19)
	obj.Data = append(obj.Data, patch20)
	obj.Data = append(obj.Data, patch21)
	obj.Data = append(obj.Data, patch22)
	obj.Data = append(obj.Data, patch23)
	obj.Data = append(obj.Data, patch24)

	return obj
}
