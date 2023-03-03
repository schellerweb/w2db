// ================================================================================================
// File	Name	: main.go
// Project		: w2db
// Author		: Holger Scheller
// Version		: 1.0.1
// Last Update	: 23.02.2023
// Description	: Backend functions for the w2ui library Version 1.5.
// ================================================================================================
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const Version = "1.0.1"

var db *sql.DB

// ================================================================================================
// Function		: init
// Description	: open SQLite database
// ================================================================================================
func init() {
	db, err := sql.Open("sqlite3", "w2db.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
}

// ================================================================================================
// Function		: main
// Description	: start webServer()
// ================================================================================================
func main() {
	webServer()
}

// ================================================================================================
// Function		: webServer
// Description	: request / response w2grid / w2form url
// ================================================================================================
func webServer() {
	// init http router
	router := mux.NewRouter()
	sPath := ""
	nPort := 3000
	router.PathPrefix(sPath + "/web/").Handler(http.StripPrefix(sPath+"/web/", http.FileServer(http.Dir("./web"))))
	router.HandleFunc(sPath+"/w2grid", w2grid)
	router.HandleFunc(sPath+"/w2form", w2form)
	log.Printf("Start %s Version: %s", filepath.Base(os.Args[0]), Version)
	log.Printf("http://127.0.0.1:%d/web/index.html", nPort)
	srv := &http.Server{
		Handler: handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router),
		Addr:    ":" + strconv.Itoa(nPort),
		// enforce timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
