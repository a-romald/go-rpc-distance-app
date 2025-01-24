package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/a-romald/go-rpc-distance-app/driver"
)

type App struct{}

var client *sql.DB

func main() {

	// create logger
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// db connection
	dsn := "user:secret@tcp(mysql:3306)/geodb?parseTime=true&tls=false"
	conn, err := driver.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := App{}

	client = conn

	// Register the RPC Server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic()
	}

	app.rpcListen()
}

func (app *App) rpcListen() error {
	log.Println("Starting RPC server on port ", os.Getenv("RPC_PORT"))
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("RPC_PORT")))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
