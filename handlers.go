package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Json file for the maps data
// https://data.sfgov.org/resource/rqzj-sfat.json
// AIzaSyCYy0Pt6UolytUOtxbdFdGkUA3iao0UkrA -- serverkey for google distances
var (
	tdb *TruckDB
)

func initializeTruckLocationServer() error {
	initializeFoodCategories()

	tdb = NewTruckDB()
	e := tdb.readJsonFile()
	if e != nil {
		return e
	}

	tdb.updateKDTree()
	return nil
}

func ErrorHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadFile("errors/index.html")
	fmt.Fprint(rw, string(body))
}

func TrucksHandler(rw http.ResponseWriter, req *http.Request) {

	kv := req.URL.Query()

	l1, _ := strconv.ParseFloat(kv["lat"][0], 64)
	l2, _ := strconv.ParseFloat(kv["lng"][0], 64)
	rad, _ := strconv.ParseFloat(kv["radius"][0], 64)
	count, _ := strconv.ParseInt(kv["count"][0], 0, 64)
	types := kv["types"][0]

	ll := LatLng{l1, l2}
	categories := []string{}
	if len(types) > 0 {
		ws := strings.Split(types, ",")
		for _, w := range ws {
			categories = append(categories, w)
		}
	}
	// THIS is to find all the locations withing a given radius
	var res TrucksRes = tdb.findNearestTrucks(ll, rad/1600.0, count, categories)
	json.NewEncoder(rw).Encode(res)
}
