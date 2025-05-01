package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	processsalesfile "main.go/ProcessSalesFile"
	revenue "main.go/Revenue"
	"main.go/db"
	"main.go/middleware"
)

func main() {
	log.Println("main (+)")

	configFile := flag.String("config", "", "Run with configuration file, refer Readme.md")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("Configuration file not provided. Use -config <dev.json>")
	}

	err := db.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	lFile, lErr := os.OpenFile("./log/log"+time.Now().Format("20060102150405")+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatal("Failed to open the file :", lErr.Error())
	}

	defer lFile.Close()

	log.SetOutput(lFile)

	lErr = db.ConnectDB()
	if lErr != nil {
		log.Println("Error Main 001 :", lErr.Error())
	}

	db.RunMigrations()

	go db.StartDailyRefreshJob()

	lRouter := mux.NewRouter()

	// lRouter.HandleFunc("/salesdata/read", processsalesfile.ReadAndStoreSalesFileData).Methods(http.MethodPut)

	// lErr = http.ListenAndServe(":28050", lRouter)

	lRouter.HandleFunc("/salesdata/read", processsalesfile.ReadAndStoreSalesFileData).Methods(http.MethodPut)
	lRouter.HandleFunc("/revenue/total", revenue.HandleTotalRevenue).Methods(http.MethodGet)
	lRouter.HandleFunc("/revenue/byProduct", revenue.HandleTotalRevenueByProduct).Methods(http.MethodPut)
	lRouter.HandleFunc("/revenue/byCategory", revenue.HandleTotalRevenueByCategory).Methods(http.MethodPut)
	lRouter.HandleFunc("/revenue/byRegion", revenue.HandleTotalRevenueByRegion).Methods(http.MethodPut)
	lRouter.HandleFunc("/refresh", db.RefreshDataHandler)

	handlerWithMiddleware := middleware.LoggingMiddleware(db.GPostgres)(lRouter)
	http.ListenAndServe(":28050", handlerWithMiddleware)

	if lErr != nil {
		log.Println("Error Main 002 :", lErr.Error())
		log.Fatal(lErr)
	}

	log.Println("main (-)")
}
