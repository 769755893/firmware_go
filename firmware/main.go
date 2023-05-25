package main

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

// "github.com/petarov/query-apple-firmware-updates/config"

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

const INFO_URL = "https://api.ipsw.me/v4/device/%s?type=ipsw"


func getResponseBody(resp *http.Response) io.ReadCloser {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	return reader
}

func decodeIpsw(data []byte) (*IPSWFirm, error) {
	var r IPSWFirm
	err := json.Unmarshal(data, &r)
	return &r, err
}

func encodeIpsw(r *IPSWFirm) ([]byte, error) {
	return json.Marshal(r)
}


func IPSWGetInfo(product string) (ipsw *IPSWFirm, err error) {
	url := fmt.Sprintf(INFO_URL, product)
	fmt.Printf("url: %v\n", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error get ")
		return nil, err
	}

	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	
	ips, err := decodeIpsw(buffer)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	return ips, nil
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

	r.Get("/getFirmware", func (w http.ResponseWriter, r *http.Request)  {
		firmwareCode := r.FormValue("firmwareCode")
		ipsw, err := IPSWGetInfo(firmwareCode)
		if err != nil {
			w.Write([]byte("error"))
		}
		jsonstr, err := encodeIpsw(ipsw)
		if err != nil {
			w.Write([] byte("error"))
		}
		w.Write([] byte(jsonstr))
	})

	http.ListenAndServe(":3333", r)
}

func mapToSignedFirmware(ipsw *IPSWFirm) (firmware [] Firmware) {
	fimwares := ipsw.Firmwares
	var ret = make([] Firmware, 1)
	for i := 0; i < len(fimwares); i++ {
		if (fimwares[i].Signed) {
			ret = append(ret, fimwares[i])
		}
	}
	return ret;
}

func saveIpswToDB(firmware [] Firmware) {

}

func main() {
	mysqlDB = ConnMySQL()
	setupRoute()
}