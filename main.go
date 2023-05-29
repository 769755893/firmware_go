package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"main/db"
	"main/route"
	"main/service"
	"sync"
)

func main() {
	db.LoadDbConfig()
	db.DeleteAllTable()
	err := db.InitTable()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		fmt.Printf("\"failed to start server\": %v\n", "failed to start server")
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Printf("\"main invoke monitor\": %v\n", "main invoke monitor")
		service.StartMonitor()
		fmt.Printf("\"main end to invoke monitor\": %v\n", "main end to invoke monitor")
	}()

	go func() {
		defer wg.Done()
		route.SetupRoute()
	}()
	wg.Wait()
}
