package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Measurement struct {
    Id			uint64	`json:"id"`
    SensorId	uint64	`json:"sensor_id"`
    Value		float64	`json:"value"`
    CreatedAt   time.Time	`json:"created_at"`
}

// Queries the database for github activities and transforms them
// to JSON format
func getMeasurementsFromDb(sensor_id uint64) ([]Measurement, error) {
	// Select all measurements
	rows, err := DB.Query("SELECT * FROM measurements WHERE sensor_id=$1 ORDER BY created_at DESC LIMIT 10000", sensor_id)
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


func getMeasurementsJson(sensor_id uint64) []byte {
	measurements, err := getMeasurementsFromDb(sensor_id)
	if err != nil {
		fmt.Println(err)
	}
	result, err := json.Marshal(map[string]interface{}{"measurements": measurements})
	return result
}
