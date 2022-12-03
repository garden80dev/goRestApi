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

const (
	GetAll        = "SELECT * FROM AllRanks WHERE omega >= 0"
	GetAllByModel = "SELECT * FROM AllRanks WHERE (car_id = '%s') AND (omega >= 0) "
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getAllMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(readRankings(""))
}

func getFerrariMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(readRankings("Ferrari"))
}

func getMaseratiMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(readRankings("maserati"))
}

func getLamborghiniMeasurements(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(readRankings("lamborghini"))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readRankings(carId string) (measurements []measurement) {
	db, err := sql.Open("sqlite3", "./ranks.db")
	check(err)
	defer db.Close()

	var query string
	if carId != "" {
		query = fmt.Sprintf(GetAllByModel, carId)
	} else {
		query = GetAll
	}

	rows, err := db.Query(query)
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

	//measurements = filterOmega((measurements))
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
	http.HandleFunc("/ferrari", getFerrariMeasurements)
	http.HandleFunc("/maserati", getMaseratiMeasurements)
	http.HandleFunc("/lamborghini", getLamborghiniMeasurements)
	log.Fatal(http.ListenAndServe(":7777", nil))
}
