package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Measurement struct {
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	QUERY_LIMIT = 10000
)

func getQueryFromParams(sensor_id uint64, interval string, date string) string {
	if interval == "" {
		return "SELECT created_at, value FROM measurements WHERE sensor_id=$1 ORDER BY created_at DESC LIMIT " + strconv.Itoa(QUERY_LIMIT)
	}
	if interval == "peek" {
		return "SELECT created_at, value FROM measurements WHERE sensor_id=$1 ORDER BY created_at DESC LIMIT 1"
	}
	date_str := ""
	if date != "" && (interval == "hour" || interval == "minute") {
		date_str = `AND (created_at + INTERVAL '12 hours') BETWEEN $3 AND ($3 + INTERVAL '24 hours')`
	}
	return `SELECT diff.period, SUM(diff.result) FROM (` +
		`SELECT ` +
		`DATE_TRUNC($2, created_at + INTERVAL '12 hours') AS period, ` +
		`value - LAG(value, 1, CAST(0 AS real)) OVER (ORDER BY id) AS result ` +
		`FROM measurements ` +
		`WHERE sensor_id=$1 ` +
		date_str +
		`) diff ` +
		`GROUP BY diff.period ` +
		`ORDER BY period DESC LIMIT ` + strconv.Itoa(QUERY_LIMIT)
}

func getQueryResultsFromParams(sensor_id uint64, interval string, date string) (*sql.Rows, error) {
	sql := getQueryFromParams(sensor_id, interval, date)
	if interval == "" {
		return DB.Query(sql, sensor_id)
	}
	if interval == "hour" || interval == "minute" {
		return DB.Query(sql, sensor_id, interval, date)
	}
	return DB.Query(sql, sensor_id, interval)
}

// Queries the database for github activities and transforms them
// to JSON format
func getMeasurementsFromDb(sensor_id uint64, interval string, date string) ([]Measurement, error) {
	rows, err := getQueryResultsFromParams(sensor_id, interval, date)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var measurements []Measurement
	for rows.Next() {
		var measurement Measurement
		err := rows.Scan(&(measurement.CreatedAt), &(measurement.Value))
		if err != nil {
			fmt.Println(err)
			continue
		}
		measurements = append(measurements, measurement)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	return measurements, err
}

func getMeasurementsJson(sensor_id uint64, interval string, date string) []byte {
	measurements, err := getMeasurementsFromDb(sensor_id, interval, date)
	if err != nil {
		fmt.Println(err)
	}
	result, err := json.Marshal(map[string]interface{}{"measurements": measurements})
	return result
}
