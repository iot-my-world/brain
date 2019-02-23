package reading

import "math"

type Reading struct {
	TimeStamp int64   `json:"timeStamp" bson:"timeStamp"`
	Id        string  `json:"id" bson:"id"`
	IMEI      string  `json:"imei" bson:"imei"`
	Raw       string  `json:"raw" bson:"raw"`
	Latitude  float32 `json:"latitude" bson:"latitude"`
	Longitude float32 `json:"longitude" bson:"longitude"`
}

const earthRadiusInKm float64 = 6378.137

func DifferenceBetween(r1, r2 *Reading) float32 {
	lat1 := r1.Latitude
	lon1 := r1.Longitude
	lat2 := r2.Latitude
	lon2 := r2.Longitude
	var dLat = float64(lat2*math.Pi/180 - lat1*math.Pi/180)
	var dLon = float64(lon2*math.Pi/180 - lon1*math.Pi/180)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(float64(lat1*math.Pi/180))*math.Cos(float64(lat2*math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = earthRadiusInKm * c
	return float32(d * 1000)
}