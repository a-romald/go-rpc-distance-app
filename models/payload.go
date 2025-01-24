package models

type Geo struct {
	Locations []Location `json:"locations"`
	IPAddress string     `json:"-"`
}

type Location struct {
	Id     string `json:"id"`
	Coords Coord  `json:"coords"`
}

type Coord struct {
	Point1 Point `json:"point1"`
	Point2 Point `json:"point2"`

	//Point1 struct {
	//	Lat float64 `json:"lat"`
	//	Lng float64 `json:"lng"`
	//}
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Result struct {
	Id       string  `json:"id"`
	Distance float64 `json:"distance"`
}
