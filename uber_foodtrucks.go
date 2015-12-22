package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kdtree"
	"net/http"
	"os"
	"strconv"
)

// Json file for the maps data
// https://data.sfgov.org/resource/rqzj-sfat.json
// AIzaSyCYy0Pt6UolytUOtxbdFdGkUA3iao0UkrA -- serverkey for google distances

var globalvar int = 2

func myHandle(rw http.ResponseWriter, req *http.Request) {
	localvar := 1

	fmt.Fprintf(rw, "Local: %d Global: %d\n", localvar, globalvar)
	globalvar += 1
}

func searchMap(rw http.ResponseWriter, req *http.Request) {
	kv := req.URL.Query()

	l1, _ := strconv.ParseFloat(kv["lat"][0], 64)
	l2, _ := strconv.ParseFloat(kv["lng"][0], 64)
	rad, _ := strconv.ParseFloat(kv["radius"][0], 64)
	count, _ := strconv.ParseInt(kv["count"][0], 0, 64)

	ll := kdtree.Latlng{l1, l2}

	// THIS is to find all the locations withing a given radius
	var res kdtree.LatlngRes = kdtree.FindNearestRadius(ll, rad/1600.0)

	// Hack .. dont like this
	var rlen int = len(res)
	var i64 int64
	i64 = int64(rlen)

	// To reverse a slice
	for left, right := 0, len(res)-1; left < right; left, right = left+1, right-1 {
		res[left], res[right] = res[right], res[left]
	}

	if i64 > count {
		json.NewEncoder(rw).Encode(res[:count])
	} else {
		json.NewEncoder(rw).Encode(res)
	}
}

func loadMap(rw http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadFile("html/index.html")
	fmt.Fprint(rw, string(body))
}

func main() {

	assets := http.StripPrefix("/", http.FileServer(http.Dir("html/")))
	http.Handle("/", assets)

	http.HandleFunc("/map", loadMap)
	http.HandleFunc("/search", searchMap)

	kdtree.BuildTree()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
	}
	fmt.Println("Listening on: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port ", port)

	//http.ListenAndServe(":8888", nil)
}
