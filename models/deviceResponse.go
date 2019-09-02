package models

/*
DeviceChannel :device channel definition
*/
type DeviceChannel struct {
	R  float64
	Gr float64
	Gb float64
	B  float64
}

/*
DeviceResponse :device response
*/
type DeviceResponse struct {
	Data map[int]DeviceChannel
}

/*
NewDeviceResponse :initialize device performance
*/
func NewDeviceResponse(data [][]string) *DeviceResponse {
	if data == nil {
		return nil
	}

	obj := new(DeviceResponse)
	obj.Data = make(map[int]DeviceChannel, 0)

	for _, row := range data {
		wavelength := StrToInt(row[0])
		channel := &DeviceChannel{
			R:  StrToFloat64(row[3]),
			Gr: StrToFloat64(row[1]),
			Gb: StrToFloat64(row[2]),
			B:  StrToFloat64(row[4]),
		}
		obj.Data[wavelength] = *channel
	}

	return obj
}

/*
SerializedData :serialize pixel data
*/
func (dr *DeviceResponse) SerializedData() map[int][]float64 {
	buffer := make(map[int][]float64, 0)
	for key, rggb := range dr.Data {
		buffer[key] = []float64{rggb.R, rggb.Gr, rggb.Gb, rggb.B}
	}

	return buffer
}
