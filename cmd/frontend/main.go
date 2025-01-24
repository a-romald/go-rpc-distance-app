package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a-romald/go-rpc-distance-app/driver"
	"github.com/a-romald/go-rpc-distance-app/models"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       models.BaseModel
}

func main() {

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// db connection
	dsn := "user:secret@tcp(mysql:3306)/geodb?parseTime=true&tls=false"
	conn, err := driver.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		DB:       models.BaseModel{DB: conn},
	}

	log.Println("Starting server on port " + os.Getenv("FRONTEND_PORT"))
	// define http server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", os.Getenv("FRONTEND_PORT")),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
