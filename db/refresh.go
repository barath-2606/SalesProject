package db

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartDailyRefreshJob() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			log.Println("Starting scheduled data refresh...")
			err := refreshData()
			if err != nil {
				log.Println("Scheduled refresh failed:", err)
			} else {
				log.Println("Scheduled refresh succeeded")
			}

			<-ticker.C
		}
	}()
}

func RefreshDataHandler(w http.ResponseWriter, r *http.Request) {
	err := refreshData()
	if err != nil {
		log.Println("Manual refresh failed:", err)
		http.Error(w, "Refresh failed", http.StatusInternalServerError)
		return
	}
	log.Println("Manual refresh succeeded")
	w.Write([]byte("Refresh done"))
}

func refreshData() error {
	// Log start
	start := time.Now()
	log.Println("Data refresh started at:", start)

	tx := GPostgres.Begin()

	// Option 1: Truncate and reinsert
	if err := tx.Exec("TRUNCATE TABLE sales_file_data").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("truncate failed: %v", err)
	}

	if err := tx.Exec("TRUNCATE TABLE product_details").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("truncate failed: %v", err)
	}

	if err := tx.Exec("TRUNCATE TABLE customer_details").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("truncate failed: %v", err)
	}

	if err := tx.Exec("TRUNCATE TABLE order_details").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("truncate failed: %v", err)
	}

	tx.Commit()
	log.Println("Data refresh completed at:", time.Now())
	return nil
}
