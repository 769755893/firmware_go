package route

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"main/model"
	"main/service"
	"main/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleCallGetFirmware(w http.ResponseWriter, r *http.Request) {
	firmwareCode := r.FormValue("firmwareCode")
	if len(firmwareCode) == 0 {
		util.HttpWriter(w, 400, util.ByteString("require params firmwareCode"))
	}
	var dbRet = service.GenerateFirmware(firmwareCode)

	fmt.Printf("len(dbRet): %v\n", len(dbRet))
	if len(dbRet) == 0 {
		ipsw, err := service.GetNewIpswFirmware(firmwareCode)
		if err != nil {
			util.HttpWriter(w, 500, util.ByteString("server Get Ipsw Firmware failed"))
			return
		}
		json, err := model.EncodeIpsw(ipsw)
		if err != nil {
			util.HttpWriter(w, 500, util.ByteString("error"))
		}
		util.HttpWriter(w, 200, json)
		return
	}

	var ret = model.IPSWFirm{
		Name:       dbRet[0].Name,
		Identifier: firmwareCode,
		Firmwares:  dbRet,
	}

	str, er := model.EncodeIpsw(&ret)
	if er != nil {
		w.WriteHeader(500)
		_, err := w.Write([]byte("server encode ipsw error"))
		if err != nil {
			return
		}
	}
	util.HttpWriter(w, 200, str)
}

func SetupRoute() {
	fmt.Printf("\"start setup route\": %v\n", "start setup route")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/getFirmware", func(w http.ResponseWriter, r *http.Request) {
		handleCallGetFirmware(w, r)
	})

	// sudo lsof -i :port
	// sudo KILL -9 PID release port
	// test : http://localhost:3333/

	var server = http.Server{
		Addr: ":8080",
	}

	go func() {
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
	}()

	osShutdown := make(chan os.Signal, 1)
	signal.Notify(osShutdown, syscall.SIGINT, syscall.SIGTERM)
	<-osShutdown

	if e := server.Close(); e != nil {
		fmt.Printf("\"HTTP server shutdown failed\": %v\n", "HTTP server shutdown failed")
	}
	fmt.Printf("\"HTTP server shutdown success\": %v\n", "HTTP server shutdown success")
}
