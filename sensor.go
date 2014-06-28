package main

import (
	"encoding/json"
	"fmt"
	"time"
	"database/sql"
	"strings"
)

type Sensor struct {
    Id			int		`json:"id"`
    Type		string	`json:"type"`
    Name		string	`json:"name"`
    CreatedAt   time.Time	`json:"created_at"`
    UpdatedAt   time.Time	`json:"updated_at"`
}

// Queries the database for sensors and transforms them
// to struct format
func getSensorsFromDb(sensor_ids []uint64) ([]Sensor, error) {
	// Select all sensors
	var rows * sql.Rows
	var err error
	if sensor_ids == nil {
		rows, err = DB.Query("SELECT * FROM sensors LIMIT 50")
	} else if len(sensor_ids) == 1 {
		rows, err = DB.Query("SELECT * FROM sensors WHERE id=$1", sensor_ids[0])
	} else {
		sql := "SELECT * FROM sensors WHERE id in (?" + strings.Repeat(",?", len(sensor_ids)-1) + ")"
		rows, err = DB.Query(sql, sensor_ids)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var sensors []Sensor
	for rows.Next() {
		var sensor Sensor
		err := rows.Scan(&(sensor.Id), &(sensor.Type), &(sensor.Name), &(sensor.CreatedAt), &(sensor.UpdatedAt))
		if err != nil {
			fmt.Println(err)
			continue
		}
		sensors = append(sensors, sensor)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	return sensors, err
}

func getSensorFromDb(sensor_id uint64) (Sensor, error) {
	sensor_ids := []uint64{sensor_id}
	sensor_ids[0] = sensor_id
	sensors, err := getSensorsFromDb(sensor_ids)
	var sensor Sensor
	if err != nil {
		return sensor, err
	}
	sensor = sensors[0]
	return sensor, nil
}

func getSensorJson(sensor_id uint64) []byte {
	sensor, err := getSensorFromDb(sensor_id)
	if err != nil {
		fmt.Println(err)
	}
	result, err := json.Marshal(map[string]interface{}{"sensors": sensor})
	return result
}

func getSensorsJson(sensor_ids []uint64) []byte {
	sensors, err := getSensorsFromDb(sensor_ids)
	if err != nil {
		fmt.Println(err)
	}
	for _, sensor := range sensors {
		fmt.Println("Got sensor " + sensor.Name)
	}
	content, err := json.Marshal(sensors)
	fmt.Println(string(content))
	if err != nil {
		return nil
	}
	var content_iface interface{}
	json.Unmarshal(content, &content_iface)
	result, err := json.Marshal(map[string]interface{}{"sensors": content_iface})
	return result
}
