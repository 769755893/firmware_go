package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var mysqlDB *sql.DB

func ConnMySQL() *sql.DB {
	driverName := "mysql"

	user := "root"
	password := "00000000"


	port := "127.0.0.1:3306"
	name := "testdb1"
	dataSourceName := user + ":" + password + "@" + "tcp" + "(" + port + ")" + "/" + name + "?charset=utf8"

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetConnMaxIdleTime(time.Minute * 10)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		fmt.Println("mysql connect failed")
		panic(err)
	}

	fmt.Println("mysql connect complete")
	return db
}

func setupRoute() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("error"))
	})

	r.Get("/getFirmware", handleCallGetFirmware)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("failed to start server with err: %v\n", err)
	} else {
		fmt.Printf("\"success start server at 3333 port\": %v\n", "success start server at 3333 port")
	}
}

func handleCallGetFirmware(w http.ResponseWriter, r *http.Request) {
	firmwareCode := r.FormValue("firmwareCode")
	var dbRet = generateFirmware(*mysqlDB, firmwareCode)

	if len(dbRet) == 0 {
		ipsw, err := getNewIpswFirmware(firmwareCode)
		if err != nil {
			w.Write([] byte ("error"))
			return
		}
		jsonstr, err := encodeIpsw(ipsw)
		if err != nil {
			w.Write([] byte("error"))
		}
		w.Write([] byte(jsonstr))
		return
	}

	var ret IPSWFirm = IPSWFirm{
		Name: dbRet[0].Name,
		Identifier: firmwareCode,
		Firmwares: dbRet,
	}

	jsonRet, er := encodeIpsw(&ret)
	if er != nil {
		w.Write([] byte("error"))
	}
	w.Write([] byte (jsonRet))
}


/// already sync db data.
func getNewIpswFirmware(firmwareCode string)(*IPSWFirm, error) {
	ipsw, err := getIpswInfo(firmwareCode)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	saveIpswsToDB(firmwareCode, ipsw.Name, mapToSignedFirmware(ipsw))
	return ipsw, nil
}

func saveIpswsToDB(firmwareCode string, iphoneName string, firmwares [] Firmware) {
	updateFirmware(*mysqlDB, firmwareCode, iphoneName, firmwares, table_firmware)
}

func startMonitor() {
	tiker := time.Tick(3600 * time.Second)

	for range tiker {
		tickerUpdateFirmware(*mysqlDB)
	}
}

func main() {
	mysqlDB = ConnMySQL()
	deleteAllTable(*mysqlDB)
	err := initTable(*mysqlDB)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fmt.Printf("\"failed to start server\": %v\n", "failed to start server")
	}
	startMonitor()
	setupRoute()
}