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
		Id: "0001"
		Timestamp: "2022-10-04T15:15:47.000Z",
		Position:  "position1",
		Temp:      "50",
		Omega:     "ooo",
		Speed:     "sss",
		Car_id:    "Ferrari",
	},
	{
		Id: "0002"
		Timestamp: "2022-11-04T15:15:47.000Z",
		Position:  "position2",
		Temp:      "80",
		Omega:     "ooo2",
		Speed:     "sss2",
		Car_id:    "Maserati",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getAllMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(remark)
}

func openDb() {
	database, _ := sql.Open("sqlite3", "./ranks.db")
    //statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS AllRanks (id INTEGER PRIMARY KEY, timestamp TEXT, press INTEGER, )")
    //statement.Exec()
    //statement, _ = database.Prepare("INSERT INTO AllRanks (timestamp, press) VALUES (?, ?)")
    //statement.Exec("Nic", "Raboy")
    rows, _ := database.Query("SELECT id, timestamp, press FROM AllRanks")
    var id int
    var timestamp string
    var press int
    for rows.Next() {
        rows.Scan(&id, &timestamp, &press)
        fmt.Println(strconv.Itoa(id) + ": " + timestamp + " " + press)
    }
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/events", getAllMeasurements).Methods("GET")
	log.Fatal(http.ListenAndServe(":7777", router))
}
