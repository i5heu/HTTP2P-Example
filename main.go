package main

import (
	"database/sql"
	"fmt"
	"http2p"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
)

var db *sql.DB
var err error

type Config struct {
	Dblogin string
}

var conf Config

func main() {

	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		fmt.Println(err)
		return
	}

	// ################ END CONFIG ###########################

	// Create an sql.DB and check for errors
	db, err = sql.Open("mysql", conf.Dblogin)
	db.SetConnMaxLifetime(time.Second * 2)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(25)

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic
	}

	if err != nil {
		panic(err.Error())
	}
	// sql.DB should be long lived "defer" closes it once this function ends
	defer db.Close()

	http2p.Init(db)

	// #################### END CREATE DATABASE #####################

	fmt.Println("START")

	http.HandleFunc("/swarm", SwarmHandler)
	http.ListenAndServe(":8084", nil)

}

func SwarmHandler(w http.ResponseWriter, r *http.Request) {
	http2p.APiHandler(w, r)
}
