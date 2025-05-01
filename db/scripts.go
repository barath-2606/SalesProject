package db

import (
	"log"

	"main.go/models"
)

func RunMigrations() error {
	lTx := GPostgres.Begin()
	if lTx.Error != nil {
		log.Println("Error starting transaction:", lTx.Error)
		return lTx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			lTx.Rollback()
			log.Println("Recovered from panic, rolled back transaction.")
		}
	}()

	if err := lTx.AutoMigrate(&models.SalesFileData{}); err != nil {
		log.Println("Error RM 001:", err)
		return err
	}
	if err := lTx.AutoMigrate(&models.CustomerDetails{}); err != nil {
		log.Println("Error RM 002:", err)
		return err
	}
	if err := lTx.AutoMigrate(&models.OrderDetails{}); err != nil {
		log.Println("Error RM 003:", err)
		return err
	}
	if err := lTx.AutoMigrate(&models.ProductDetails{}); err != nil {
		log.Println("Error RM 004:", err)
		return err
	}

	if err := lTx.Commit().Error; err != nil {
		log.Println("Error RM 005 (commit):", err)
		return err
	}

	viewSQL := `
CREATE OR REPLACE VIEW daily_sales_summary AS
SELECT
  userid,
  DATE(created_at) AS sales_date,
  SUM(net) AS total_sales
FROM sales_file_data
GROUP BY userid, DATE(created_at);
`
	if err := GPostgres.Exec(viewSQL).Error; err != nil {
		log.Println("Error creating view:", err)
		return err
	}

	log.Println("Tables and view created successfully.")
	return nil
}
