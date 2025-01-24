package utils

import (
	"math"
)

const radius = 6371 // Approximate radius of earth in km

/*
   Valid latitudes are +90 to -90 degrees.
   Valid longitudes are +180 to -180 degrees but it is an angle so +/-180 are the same place and any value can represent a valid longitude.
*/

func Calculate_km_distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	if lat1 < -90 || lat1 > 90 || lat2 < -90 || lat2 > 90 {
		return 0
	}
	if lng1 < -360 || lng1 > 360 || lng2 < -360 || lng2 > 360 {
		return 0
	}
	dlat := degrees2radians(lat2 - lat1)
	dlon := degrees2radians(lng2 - lng1)
	a := (math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(degrees2radians(lat1))*
		math.Cos(degrees2radians(lat2))*math.Sin(dlon/2)*math.Sin(dlon/2))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := radius * c

	return roundFloat(d, 3)
}

// same as math.radians in python
func degrees2radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
