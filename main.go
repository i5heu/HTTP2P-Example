package main

import (
	"database/sql"
	"fmt"
	"gosocial"
	"http2p"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Config struct {
	Dblogin string
}

var conf Config

var fs = http.FileServer(http.Dir("static"))

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
	gosocial.Init(db)

	// #################### END CREATE DATABASE #####################

	fmt.Println("START")

	http.HandleFunc("/swarm", SwarmHandler)
	http.HandleFunc("/gosocial", GoSocial)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8084", nil)

}

func SwarmHandler(w http.ResponseWriter, r *http.Request) {
	http2p.APiHandler(w, r)
}

func GoSocial(w http.ResponseWriter, r *http.Request) {
	gosocial.ApiHandler(w, r, "asd")
}
