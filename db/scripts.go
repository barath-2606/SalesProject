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

	var viewSQL = []string{
		`CREATE OR REPLACE VIEW product_summary AS
	SELECT
	  product_id,
	  product_name,
	  category,
	  unit_price,
	  created_date,
	  created_by
	FROM public.product_details`,

		`CREATE OR REPLACE VIEW order_summary AS
	SELECT
	  order_id,
	  customer_id,
	  product_id,
	  quantity_sold,
	  unit_price,
	  discount,
	  shipping_cost,
	  (quantity_sold * unit_price - discount + shipping_cost) AS total_amount,
	  payment_method,
	  date_of_sale,
	  region,
	  created_by,
	  created_date
	FROM public.order_details;`,

		`CREATE OR REPLACE VIEW customer_summary AS
	SELECT
	  customer_id,
	  customer_name,
	  customer_email,
	  customer_address,
	  created_date,
	  created_by
	FROM public.customer_details`,
	}

	for _, lQuery := range viewSQL {
		if err := GPostgres.Exec(lQuery).Error; err != nil {
			log.Println("Error creating view:", err)
			return err
		}
	}

	log.Println("Tables and view created successfully.")
	return nil
}
