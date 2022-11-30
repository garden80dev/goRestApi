package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type measurement struct {
	Timestamp string `json:"Timestamp"`
	Position  string `json:"Position"`
	Temp      string `json:"Temp"`
	Omega     string `json:"Omega"`
	Speed     string `json:"Speed"`
	Car_id    string `json:"Car_id"`
}

type measurements []measurement

var remark = measurements{
	{
		Timestamp: "2022-10-04T15:15:47.000Z",
		Position:  "position1",
		Temp:      "50",
		Omega:     "ooo",
		Speed:     "sss",
		Car_id:    "0001",
	},
	{
		Timestamp: "2022-11-04T15:15:47.000Z",
		Position:  "position2",
		Temp:      "80",
		Omega:     "ooo2",
		Speed:     "sss2",
		Car_id:    "0002",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getAllMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(remark)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/events", getAllMeasurements).Methods("GET")
	log.Fatal(http.ListenAndServe(":7777", router))
}
