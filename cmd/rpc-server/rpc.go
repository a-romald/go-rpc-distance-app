package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-romald/go-rpc-distance-app/models"
	"github.com/a-romald/go-rpc-distance-app/utils"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available over RPC.
type RPCServer struct{}

func (r *RPCServer) DistanceInfo(geo models.Geo, resp *string) error {

	var results []*models.Result
	var geoLocs []*models.GeoLocation
	c1 := make(chan float64, 1)
	mod := models.BaseModel{DB: client}

	for _, data := range geo.Locations {
		var dist float64
		var res models.Result
		var geoLoc models.GeoLocation

		point1 := data.Coords.Point1
		point2 := data.Coords.Point2
		go func() {
			dist = utils.Calculate_km_distance(point1.Lat, point1.Lng, point2.Lat, point2.Lng)
			c1 <- dist
		}()
		res.Id = data.Id
		res.Distance = <-c1
		results = append(results, &res)
		// new geoLoc to database
		geoLoc.Point1.Lat = point1.Lat
		geoLoc.Point1.Lng = point1.Lng
		geoLoc.Point2.Lat = point2.Lat
		geoLoc.Point2.Lng = point2.Lng
		geoLoc.Distance = res.Distance
		geoLoc.IpAddress = geo.IPAddress
		geoLoc.CreatedAt = time.Now()
		geoLocs = append(geoLocs, &geoLoc)
	}
	// close the channel
	close(c1)

	jsonObj, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
	}

	// Insert into DB
	go func() {
		for _, item := range geoLocs {
			if item.Distance > 0 {
				_, err := mod.GetLocation(*item) // err if locations pair not found
				if err != nil {
					mod.InsertLocation(*item)
				}
			}
		}
	}()

	// resp is the message sent back to the RPC caller
	*resp = string(jsonObj)
	return nil
}
