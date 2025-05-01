package revenue

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"main.go/db"
	"main.go/models"
)

func HandleTotalRevenue(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleTotalRevenue (+)")
	startStr := r.Header.Get("start")
	endStr := r.Header.Get("end")

	var lResp models.RevenueRespStruct
	var lErr error
	lResp.Status = "S"

	lResp.TotRev, lErr = getTotRev(startStr, endStr)
	if lErr != nil {
		log.Println("Error HTR 001 :", lErr.Error())
		lResp.Status = "E"
		lResp.ErrMsg = "Error HTR 001 :" + lErr.Error()
	}

	lData, lErr := json.Marshal(lResp)
	if lErr != nil {
		log.Println("Error HTR 002 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}
	log.Println("HandleTotalRevenue (-)")
}

func getTotRev(lStartDate, lEndDate string) (float64, error) {
	log.Println("getTotRev (+)")

	lTotRevenue := 0.0

	lErr := db.GPostgres.Table("order_details").Where("date_of_sale <= ? and date_of_sale >= ?", lStartDate, lEndDate).Select("COALESCE(SUM((unit_price * quantity_sold) * (1 - discount)), 0) AS total_sales").Scan(&lTotRevenue).Error

	if lErr != nil {
		log.Println("Error gTR 001 :", lErr.Error())
		return lTotRevenue, lErr
	}

	log.Println("getTotRev (-)")
	return lTotRevenue, nil
}

func HandleTotalRevenueByProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleTotalRevenueByProduct (+)")

	if r.Method != http.MethodGet {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	var lreq models.ProductRevReq
	var lRevResp models.RevenueRespStruct
	lRevResp.Status = "S"
	var lErr error

	lBody, lErr := io.ReadAll(r.Body)
	if lErr != nil {
		log.Println("Error HTRBP 001 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	lErr = json.Unmarshal(lBody, &lreq)
	if lErr != nil {
		log.Println("Error HTRBP 002 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	if lreq.ProductName == "" {
		log.Println("Error HTRBP 003 : Product Name should not be empty")
		lRevResp.Status = "E"
		lRevResp.ErrMsg = "Product Name should not be empty"
	}

	lRevResp.TotRev, lErr = getByProduct(lreq.StartDate, lreq.EndDate, lreq.ProductName)
	if lErr != nil {
		log.Println("Error HTRBP 004 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	lData, lErr := json.Marshal(lRevResp)
	if lErr != nil {
		log.Println("Error HTRBP 005 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("HandleTotalRevenueByProduct (-)")
}

func getByProduct(str, end, product string) (float64, error) {
	log.Println("getByProduct (+)")

	var lProductRevenue float64
	var lProductId string

	lErr := db.GPostgres.Table("product_details").Where("product_name = ?", product).Select("product_id").Scan(&lProductId).Error
	if lErr != nil {
		log.Println("Error gBP 001 :", lErr.Error())
		return lProductRevenue, lErr
	}

	lErr = db.GPostgres.Table("order_details").Where("product_id = ? and date_of_sale <= ? and date_of_sale >= ?", lProductId, str, end).Select("coalesce(SUM((unit_price * quantity_sold) * (1 - discount)),0) AS total_sales").Scan(&lProductRevenue).Error
	if lErr != nil {
		log.Println("Error gBP 002 :", lErr.Error())
		return lProductRevenue, lErr
	}

	log.Println("getByProduct (-)")
	return lProductRevenue, nil
}

func HandleTotalRevenueByCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleTotalRevenueByProduct (+)")

	if r.Method != http.MethodGet {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	var lreq models.CategoryRevReq
	var lRevResp models.RevenueRespStruct
	lRevResp.Status = "S"
	var lErr error

	lBody, lErr := io.ReadAll(r.Body)
	if lErr != nil {
		log.Println("Error HTRBP 001 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	lErr = json.Unmarshal(lBody, &lreq)
	if lErr != nil {
		log.Println("Error HTRBP 002 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	if lreq.Category == "" {
		log.Println("Error HTRBP 003 : Category should not be empty")
		lRevResp.Status = "E"
		lRevResp.ErrMsg = "Category should not be empty"
	}

	lRevResp.TotRev, lErr = getBycategory(lreq.StartDate, lreq.EndDate, lreq.Category)
	if lErr != nil {
		log.Println("Error HTRBP 004 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	lData, lErr := json.Marshal(lRevResp)
	if lErr != nil {
		log.Println("Error HTRBP 005 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("HandleTotalRevenueByProduct (-)")
}

func getBycategory(str, end, category string) (float64, error) {
	log.Println("getByProduct (+)")

	var lCategoryRevenue float64
	var lProductId string

	lErr := db.GPostgres.Table("product_details").Where("category = ?", category).Select("product_id").Scan(&lProductId).Error
	if lErr != nil {
		log.Println("Error gBP 001 :", lErr.Error())
		return lCategoryRevenue, lErr
	}

	lErr = db.GPostgres.Table("order_details").Where("product_id = ? and date_of_sale <= ? and date_of_sale >= ?", lProductId, str, end).Select("coalesce(SUM((unit_price * quantity_sold) * (1 - discount)),0) AS total_sales").Scan(&lCategoryRevenue).Error
	if lErr != nil {
		log.Println("Error gBP 002 :", lErr.Error())
		return lCategoryRevenue, lErr
	}

	log.Println("getByProduct (-)")
	return lCategoryRevenue, nil
}

func HandleTotalRevenueByRegion(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleTotalRevenueByProduct (+)")

	if r.Method != http.MethodGet {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	var lreq models.RegionRevReq
	var lRevResp models.RevenueRespStruct
	lRevResp.Status = "S"
	var lErr error

	lBody, lErr := io.ReadAll(r.Body)
	if lErr != nil {
		log.Println("Error HTRBP 001 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	lErr = json.Unmarshal(lBody, &lreq)
	if lErr != nil {
		log.Println("Error HTRBP 002 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	if lreq.Region == "" {
		log.Println("Error HTRBP 003 : Region should not be empty")
		lRevResp.Status = "E"
		lRevResp.ErrMsg = "Category should not be empty"
	}

	lRevResp.TotRev, lErr = getByRegion(lreq.StartDate, lreq.EndDate, lreq.Region)
	if lErr != nil {
		log.Println("Error HTRBP 004 :", lErr.Error())
		lRevResp.Status = "E"
		lRevResp.ErrMsg = lErr.Error()
	}

	lData, lErr := json.Marshal(lRevResp)
	if lErr != nil {
		log.Println("Error HTRBP 005 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("HandleTotalRevenueByProduct (-)")
}

func getByRegion(str, end, region string) (float64, error) {
	log.Println("getByProduct (+)")

	var lRegionRevenue float64

	lErr := db.GPostgres.Table("order_details").Where("region = ? and date_of_sale <= ? and date_of_sale >= ?", region, str, end).Select("coalesce(SUM((unit_price * quantity_sold) * (1 - discount)),0) AS total_sales").Scan(&lRegionRevenue).Error
	if lErr != nil {
		log.Println("Error gBP 002 :", lErr.Error())
		return lRegionRevenue, lErr
	}

	log.Println("getByProduct (-)")
	return lRegionRevenue, nil
}
