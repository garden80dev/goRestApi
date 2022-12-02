package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type measurement struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Position  string `json:"position"`
	Press     string `json:"press"`
	Temp      string `json:"temp"`
	Omega     int    `json:"omega"`
	Speed     string `json:"speed"`
	CarID     string `json:"car_id"`
}

var remark = []measurement{
	{
		ID:        "0001",
		Timestamp: "2022-10-04T15:15:47.000Z",
		Position:  "position1",
		Temp:      "50",
		Omega:     420,
		Speed:     "sss",
		CarID:     "Ferrari",
		Press:     "mamt",
	},
	{
		ID:        "0002",
		Timestamp: "2022-11-04T15:15:47.000Z",
		Position:  "position2",
		Temp:      "80",
		Omega:     -77,
		Speed:     "sss2",
		CarID:     "Maserati",
		Press:     "mamt",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getAllMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(readRankings())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readRankings() (measurements []measurement) {
	db, err := sql.Open("sqlite3", "./ranks.db")
	check(err)

	rows, err := db.Query("SELECT * FROM AllRanks")
	check(err)

	for rows.Next() {
		var m measurement

		err := rows.Scan(
			&m.ID,
			&m.Timestamp,
			&m.Press,
			&m.Position,
			&m.Temp,
			&m.Omega,
			&m.Speed,
			&m.CarID,
		)
		check(err)

		measurements = append(measurements, m)
	}

	measurements = filterOmega((measurements))
	sort.Slice(measurements, func(i, j int) bool {
		m1t, err := time.Parse(time.RFC3339, measurements[i].Timestamp)
		check(err)

		m2t, err := time.Parse(time.RFC3339, measurements[j].Timestamp)
		check(err)

		return m1t.Before(m2t)
	})

	return
}

func filterOmega(m []measurement) []measurement {
	var filtered []measurement

	for _, me := range m {
		if me.Omega >= 0 {
			filtered = append(filtered, me)
		}
	}

	return filtered
}

func main() {
	http.HandleFunc("/", homeLink)
	http.HandleFunc("/allranks", getAllMeasurements)
	log.Fatal(http.ListenAndServe(":7777", nil))
}
