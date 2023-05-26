package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var mysqlDB *sql.DB

func openDB() *sql.DB {
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

func dbWell() bool {
	err := mysqlDB.Ping()
	if err != nil {
		return false
	}
	return true
}

func setupRoute() {
	fmt.Printf("\"start setup route\": %v\n", "start setup route")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("error"))
	})

	r.Get("/getFirmware", func (w http.ResponseWriter, r *http.Request)  {
		mysqlDB = openDB()
		handleCallGetFirmware(w,r)
		defer mysqlDB.Close()	
	})

	// sudo KILL -9 PID release port
	// test : http://localhost:3333/

	var server = http.Server {
		Addr: ":8080",
	}

	go func ()  {
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}()

	osShutdown := make(chan os.Signal,1)
	signal.Notify(osShutdown, syscall.SIGINT, syscall.SIGTERM)
	<-osShutdown

	if e := server.Close(); e != nil {
		fmt.Printf("\"HTTP server shutdown failed\": %v\n", "HTTP server shutdown failed")
	}
	fmt.Printf("\"HTTP server shutdown success\": %v\n", "HTTP server shutdown success")
}

func handleCallGetFirmware(w http.ResponseWriter, r *http.Request) {
	firmwareCode := r.FormValue("firmwareCode")
	var dbRet = generateFirmware(*mysqlDB, firmwareCode)

	fmt.Printf("len(dbRet): %v\n", len(dbRet))
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
	fmt.Printf("\"getNewIpswFirmware\": %v\n", "getNewIpswFirmware")
	ipsw, err := getIpswInfo(firmwareCode)
	fmt.Printf("get the new ipsw: %v\n", ipsw)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	saveIpswsToDB(firmwareCode, ipsw.Name, mapToSignedFirmware(ipsw))
	return mapToSignedIpsw(ipsw), nil
}

func saveIpswsToDB(firmwareCode string, iphoneName string, firmwares [] Firmware) {
	fmt.Printf("\"saveIpswToDB\": %v\n", "saveIpswToDB")
	fmt.Printf("firmwares: %v\n", firmwares)
	updateFirmware(*mysqlDB, firmwareCode, iphoneName, firmwares, table_firmware)
}

func startMonitor() {
	tiker := time.Tick(3600 * time.Second)
	fmt.Printf("\"monitor started\": %v\n", "monitor started")

	for range tiker {
		fmt.Printf("\"monitor running update the firmware\": %v\n", "monitor running update the firmware")
		tickerUpdateFirmware(*mysqlDB)
	}
}

func main() {
	mysqlDB = openDB()
	deleteAllTable(*mysqlDB)
	err := initTable(*mysqlDB)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fmt.Printf("\"failed to start server\": %v\n", "failed to start server")
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Printf("\"main invoke monitor\": %v\n", "main invoke monitor")
		startMonitor()
		fmt.Printf("\"main end to invoke monitor\": %v\n", "main end to invoke monitor")
	}()

	go func ()  {
		defer wg.Done()
		setupRoute()
	}()
	wg.Wait()
}