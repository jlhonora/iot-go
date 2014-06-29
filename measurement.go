package main

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"database/sql"
)

type Measurement struct {
    Id			uint64	`json:"id"`
    SensorId	uint64	`json:"sensor_id"`
    Value		float64	`json:"value"`
    CreatedAt   time.Time	`json:"created_at"`
}

const (
	QUERY_LIMIT = 10000
)

func getQueryFromParams(sensor_id uint64, interval string) string {
	if interval == "" {
		return "SELECT * FROM measurements WHERE sensor_id=$1 ORDER BY created_at DESC LIMIT " + strconv.Itoa(QUERY_LIMIT)
	}
	if interval == "peek" {
		return "SELECT * FROM measurements WHERE sensor_id=$1 ORDER BY created_at DESC LIMIT 1"
	}
	return `SELECT diff.period, sum(diff.result) FROM (` +
			`SELECT ` +
			`date_trunc($2, interm.created_at + interval '6 hours') AS period, ` +
			`interm.value - lag(interm.value, 1, CAST(0 AS real)) OVER (ORDER BY interm.id) AS result ` +
			`FROM ` +
			`(` +
				`SELECT * ` +
				`FROM measurements ` +
				`WHERE sensor_id=$1 ` +
			`) interm ` +
			`) diff ` +
			`GROUP BY diff.period ` +
			`ORDER BY period LIMIT ` + strconv.Itoa(QUERY_LIMIT)
}

func getQueryResultsFromParams(sensor_id uint64, interval string) (*sql.Rows, error) {
	sql := getQueryFromParams(sensor_id, interval)
	fmt.Println(sql)
	if interval == "" {
		return DB.Query(sql, sensor_id)
	}
	return DB.Query(sql, sensor_id, interval)
}

// Queries the database for github activities and transforms them
// to JSON format
func getMeasurementsFromDb(sensor_id uint64, interval string) ([]Measurement, error) {
	rows, err := getQueryResultsFromParams(sensor_id, interval)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var measurements []Measurement
	for rows.Next() {
		var measurement Measurement
		err := rows.Scan(&(measurement.Id), &(measurement.SensorId), &(measurement.Value), &(measurement.CreatedAt))
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


func getMeasurementsJson(sensor_id uint64, interval string) []byte {
	measurements, err := getMeasurementsFromDb(sensor_id, interval)
	if err != nil {
		fmt.Println(err)
	}
	result, err := json.Marshal(map[string]interface{}{"measurements": measurements})
	return result
}
