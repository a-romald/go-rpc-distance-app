package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet/v6"

	"github.com/a-romald/go-rpc-distance-app/models"
	"github.com/a-romald/go-rpc-distance-app/utils"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	data := make(jet.VarMap)

	data.Set("title", "Distance Calculation Application")

	err := Render(w, "home.jet", data)
	if err != nil {
		_, _ = fmt.Fprint(w, "Error executing template:", err)
	}
}

func (app *Application) PostHandler(w http.ResponseWriter, r *http.Request) {
	/* Example:
	var str = `{"locations":[{"id":"0","coords":{"point1":{"lat":59.8730405,"lng":30.3790174},"point2":{"lat":59.8730407,"lng":30.3790176}}},{"id":"1","coords":{"point1":{"lat":59.8730407,"lng":30.3790176},"point2":{"lat":59.8730409,"lng":30.3790178}}},{"id":"3","coords":{"point1":{"lat":60.0178148,"lng":30.3872173},"point2":{"lat":60.0178146,"lng":30.3872175}}}]}`*/

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		defer r.Body.Close()

		geo := models.Geo{}
		json_err := json.Unmarshal([]byte(body), &geo)
		if json_err != nil {
			fmt.Println(json_err)
		}
		//fmt.Println(geo.Locations[0].Coords.Point1.Lat)
		// Get client IP address
		geo.IPAddress = string(r.RemoteAddr)

		// Connect to RPC-Server and transfer data
		client, err := rpc.Dial("tcp", fmt.Sprintf("rpc-server:%s", os.Getenv("RPC_PORT")))
		if err != nil {
			fmt.Println(err)
			return
		}

		var result string
		err = client.Call("RPCServer.DistanceInfo", geo, &result)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Return result
		utils.PrintJSON(w, http.StatusOK, result, "result")

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *Application) Results(w http.ResponseWriter, r *http.Request) {
	data := make(jet.VarMap)
	var sort_by string
	var sort_class string

	data.Set("title", "Distance Calculation Results")

	// Pagination
	per_page, err := strconv.Atoi(os.Getenv("PAGE_SIZE"))
	if err != nil {
		per_page = 50
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	// Sorting
	sort_by = r.URL.Query().Get("sort")

	if strings.HasPrefix(sort_by, "-") {
		sort_class = "asc"
	} else {
		sort_class = "desc"
	}
	data.Set("sort_class", sort_class)

	sort_map := make(map[string]string)
	switch sort_by {
	case "id":
		sort_map["ID"] = "-id"
		sort_map["Distance"] = "distance"
		sort_map["IPAddress"] = "ip_address"
		sort_map["CreatedAt"] = "created_at"
	case "-id":
		sort_map["ID"] = "id"
		sort_map["Distance"] = "distance"
		sort_map["IPAddress"] = "ip_address"
		sort_map["CreatedAt"] = "created_at"
	case "distance":
		sort_map["Distance"] = "-distance"
		sort_map["ID"] = "id"
		sort_map["IPAddress"] = "ip_address"
		sort_map["CreatedAt"] = "created_at"
	case "-distance":
		sort_map["Distance"] = "distance"
		sort_map["ID"] = "id"
		sort_map["IPAddress"] = "ip_address"
		sort_map["CreatedAt"] = "created_at"
	case "ip_address":
		sort_map["IPAddress"] = "-ip_address"
		sort_map["ID"] = "id"
		sort_map["Distance"] = "distance"
		sort_map["CreatedAt"] = "created_at"
	case "-ip_address":
		sort_map["IPAddress"] = "ip_address"
		sort_map["ID"] = "id"
		sort_map["Distance"] = "distance"
		sort_map["CreatedAt"] = "created_at"
	case "created_at":
		sort_map["CreatedAt"] = "-created_at"
		sort_map["ID"] = "id"
		sort_map["Distance"] = "distance"
		sort_map["IPAddress"] = "ip_address"
	case "-created_at":
		sort_map["CreatedAt"] = "created_at"
		sort_map["ID"] = "id"
		sort_map["Distance"] = "distance"
		sort_map["IPAddress"] = "ip_address"
	default:
		sort_map["ID"] = "id"
		sort_map["Distance"] = "distance"
		sort_map["IPAddress"] = "ip_address"
		sort_map["CreatedAt"] = "created_at"
	}
	data.Set("sort_map", sort_map)

	// Pagination numbers
	pages := app.DB.Paginate("location", per_page, page, sort_by)
	data.Set("pages", pages)

	// Data from db
	locations, err := app.DB.GetAllLocations(per_page, page, sort_by)
	if err != nil {
		fmt.Println(err)
	}
	data.Set("locations", locations)

	err = Render(w, "results.jet", data)
	if err != nil {
		_, _ = fmt.Fprint(w, "Error executing template:", err)
	}
}
