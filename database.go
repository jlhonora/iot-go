package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

type DBconfig struct {
	Name     string
	User     string
	Password string
}

func getDBconfig() (string, error) {
	file, err := os.Open("dbconf.json")
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(file)
	config := DBconfig{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
	return "host=localhost dbname=" + config.Name + " user=" + config.User + " sslmode=disable", err
}

func dbinit() error {
	db_str, err := getDBconfig()
	if err != nil {
		fmt.Println("Couldn't read DB config file, exiting")
	}
	fmt.Println("Using config: " + db_str)
	db, err := sql.Open("postgres", db_str)
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
	} else {
		log.Println("DB opened")
		DB = db
	}
	err = db.Ping()
	if err != nil {
		log.Println("Ping failed")
	} else {
		log.Println("Ping success")
	}
	return err
}

func dbclose() error {
	DB.Close()
	return nil
}

func getPassword(username string) string {
	var password string
	log.Printf("Querying password for " + username)
	err := DB.QueryRow("SELECT password FROM users WHERE name = $1", username).Scan(&password)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Password is %s\n", password)
	}
	return password
}
