package main

import (
	"encoding/json"
	"io"
	"fmt"
	"time"
)

type Sensor struct {
    Id			int		`json:"id"`
    Type		string	`json:"type"`
    Name		string	`json:"name"`
    CreatedAt   time.Time	`json:"created_at"`
    UpdatedAt   time.Time	`json:"updated_at"`
}

func DecodeSensor(r io.Reader) (x *Sensor, err error) {
    x = new(Sensor)
    err = json.NewDecoder(r).Decode(x)
    return
}

// Queries the database for github activities and transforms them
// to JSON format
func getSensorsFromDb() ([]Sensor, error) {
	// Select all sensors
	rows, err := DB.Query("SELECT * FROM sensors LIMIT 50")
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


func getSensorsJson() []byte {
	sensors, err := getSensorsFromDb()
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
