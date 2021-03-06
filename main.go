package main

import (
	"fmt"
	"net/http"
	"os"

	"uber_foodtruck_locator/Godeps/_workspace/src/github.com/gorilla/mux"
)

// Json file for the maps data
// https://data.sfgov.org/resource/rqzj-sfat.json

func main() {
	r := mux.NewRouter()
	e := initializeTruckLocationServer()

	r.HandleFunc("/Trucks", TrucksHandler)

	if e != nil {
		r.HandleFunc("/", ErrorHandler)
	} else {
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("html/"))))
	}

	http.Handle("/", r)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
	}
	fmt.Println("Listening on port ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
