package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
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
		log.Fatal(err)
	} else {
		log.Println("DB opened")
		DB = db
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping failed")
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
	username = strings.TrimSpace(username)
	log.Printf("Querying password for (%v)" + username, len(username))
	err := DB.QueryRow("SELECT key FROM users WHERE name = $1", username).Scan(&password)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
		break
	case err != nil:
		log.Fatal(err)
		break
	default:
		fmt.Printf("Password is %s\n", password)
		break
	}
	return password
}

func checkKey(key string, method string) bool {
	log.Printf("Querying key %v (%v) with method %v", key, len(key), method)
	key = strings.TrimSpace(key)
	var name string
	err := DB.QueryRow("SELECT name FROM users WHERE key = $1", key).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Not found")
		return false
	case err != nil:
		log.Println(err)
		return false
	default:
		log.Printf("OK")
		return true
	}
}
