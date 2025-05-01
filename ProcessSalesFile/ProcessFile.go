package processsalesfile

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"main.go/db"
	"main.go/models"
)

func ReadAndStoreSalesFileData(w http.ResponseWriter, r *http.Request) {
	log.Println("ReadAndStoreSalesFileData (+)")

	if !(r.Method == http.MethodPost) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	var lOrderDetailsArr []models.OrderDetails
	var lProductDetailsArr []models.ProductDetails
	var lCustomerDetailsArr []models.CustomerDetails

	var lWg = &sync.WaitGroup{}

	// var lInsertedIds string
	var lResp models.ReadDataResp
	lResp.Status = "S"

	lErr := r.ParseMultipartForm(10 << 20)
	if lErr != nil {
		log.Println("Error RASSD 001 :", lErr.Error())
		lResp.Status = "E"
		lResp.ErrMsg = "Error RASSD 001 :" + lErr.Error()
	}

	lFiles := r.MultipartForm.File["files"]

	lOrderDetailsArr, lProductDetailsArr, lCustomerDetailsArr, lErr = ProcessSalesDataFile(lFiles)
	if lErr != nil {
		log.Println("Error RASSD 002 :", lErr.Error())
		lResp.Status = "E"
		lResp.ErrMsg = "Error RASSD 002 :" + lErr.Error()
	}

	if len(lOrderDetailsArr) > 0 {
		lWg.Add(1)
		go func() {
			defer lWg.Done()
			lErr = ProcessRequest(lOrderDetailsArr)
			if lErr != nil {
				log.Println("Error RASSD 003 :", lErr.Error())
				lResp.Status = "E"
				lResp.ErrMsg = "Error RASSD 003 :" + lErr.Error()
			}
		}()
	}

	if len(lProductDetailsArr) > 0 {
		lWg.Add(1)
		go func() {
			defer lWg.Done()
			lErr = ProcessRequest(lProductDetailsArr)
			if lErr != nil {
				log.Println("Error RASSD 004 :", lErr.Error())
				lResp.Status = "E"
				lResp.ErrMsg = "Error RASSD 004 :" + lErr.Error()
			}
		}()
	}

	if len(lCustomerDetailsArr) > 0 {
		lWg.Add(1)
		go func() {
			defer lWg.Done()
			lErr = ProcessRequest(lCustomerDetailsArr)
			if lErr != nil {
				log.Println("Error RASSD 005 :", lErr.Error())
				lResp.Status = "E"
				lResp.ErrMsg = "Error RASSD 005 :" + lErr.Error()
			}
		}()
	}

	lWg.Wait()

	lData, lErr := json.Marshal(lResp)
	if lErr != nil {
		log.Println("Error RASSD 003 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("ReadAndStoreSalesFileData (-)")
}

func ProcessSalesDataFile(pFiles []*multipart.FileHeader) ([]models.OrderDetails, []models.ProductDetails, []models.CustomerDetails, error) {
	log.Println("ProcessSalesDataFile (+)")

	// var lInsertedIds string

	var lCustomerDet []models.CustomerDetails
	var lProductDet []models.ProductDetails
	var lOrderDet []models.OrderDetails

	for _, lFileHeder := range pFiles {
		lFile, lErr := lFileHeder.Open()
		if lErr != nil {
			log.Println("Error PSDF 001 :", lErr.Error())
			return lOrderDet, lProductDet, lCustomerDet, lErr
		}
		defer lFile.Close()

		lReader := csv.NewReader(bufio.NewReader(lFile))

		lHeader, lErr := lReader.Read()
		if lErr != nil {
			log.Println("Error PSDF 002 :", lErr.Error())
			return lOrderDet, lProductDet, lCustomerDet, lErr
		}

		lColIndex := make(map[string]int)
		for i, lVal := range lHeader {
			lColIndex[lVal] = i
		}

		lRecords, lErr := lReader.ReadAll()
		if lErr != nil {
			log.Println("Error PSDF 003 :", lErr.Error())
			return lOrderDet, lProductDet, lCustomerDet, lErr
		}

		lOrderArr, lProductArr, lCustomerArr, lErr := InsertFileRecords(lRecords, lColIndex)
		if lErr != nil {
			log.Println("Error PSDF 004 :", lErr.Error())
			return lOrderDet, lProductDet, lCustomerDet, lErr
		}

		if len(lOrderArr) > 0 {
			lOrderDet = append(lOrderDet, lOrderArr...)
		}
		if len(lProductArr) > 0 {
			lProductDet = append(lProductDet, lProductArr...)
		}
		if len(lCustomerArr) > 0 {
			lCustomerDet = append(lCustomerDet, lCustomerArr...)
		}

	}

	log.Println("ProcessSalesDataFile (-)")
	return lOrderDet, lProductDet, lCustomerDet, nil
}

func InsertFileRecords(pRecords [][]string, pColIndex map[string]int) ([]models.OrderDetails, []models.ProductDetails, []models.CustomerDetails, error) {
	log.Println("InsertFileRecords (+)")

	var lFileDataArr []models.SalesFileData
	var lCustomerArr []models.CustomerDetails
	var lOrderArr []models.OrderDetails
	var lProductArr []models.ProductDetails
	var lErr error

	lUniqueId := time.Now().Format("20060102150405.000")
	lUniqueId = strings.Replace(lUniqueId, ".", "", 1)

	lTx := db.GPostgres.Begin()
	defer lTx.Rollback()

	for _, lRecord := range pRecords {

		var lFileData models.SalesFileData

		lFileData.OrderId = cleanString(lRecord[pColIndex["Order ID"]])
		lFileData.ProductId = cleanString(lRecord[pColIndex["Product ID"]])
		lFileData.CustomerId = cleanString(lRecord[pColIndex["Customer ID"]])
		lFileData.ProductName = cleanString(lRecord[pColIndex["Product Name"]])
		lFileData.Category = cleanString(lRecord[pColIndex["Category"]])
		lFileData.Region = cleanString(lRecord[pColIndex["Region"]])
		parsedDate, _ := time.Parse("1/2/2006", lRecord[pColIndex["Date of Sale"]])
		lFileData.DateOfSale = parsedDate.Format("2006-01-02")
		lFileData.QuantitySold, _ = strconv.Atoi(lRecord[pColIndex["Quantity Sold"]])
		lFileData.UnitPrice, _ = strconv.ParseFloat(lRecord[pColIndex["Unit Price"]], 64)
		lFileData.Discount, _ = strconv.ParseFloat(lRecord[pColIndex["Discount"]], 64)
		lFileData.ShippingCost, _ = strconv.ParseFloat(lRecord[pColIndex["Shipping Cost"]], 64)
		lFileData.PaymentMethod = cleanString(lRecord[pColIndex["Payment Method"]])
		lFileData.CustomerName = cleanString(lRecord[pColIndex["Customer Name"]])
		lFileData.CustomerEmail = cleanString(lRecord[pColIndex["Customer Email"]])
		lFileData.CustomerAddress = cleanString(lRecord[pColIndex["Customer Address"]])
		lFileData.UniqueId = lUniqueId
		lFileData.CreatedBy = "AutoBot"
		lFileData.CreatedDate = time.Now().Format("2006-01-02 15:04:05")

		lCustomerRec := models.CustomerDetails{
			CustomerId:      lFileData.CustomerId,
			CustomerName:    lFileData.CustomerName,
			CustomerEmail:   lFileData.CustomerEmail,
			CustomerAddress: lFileData.CustomerAddress,
			CreatedDate:     lFileData.CreatedDate,
			CreatedBy:       lFileData.CreatedBy,
		}

		lProductRec := models.ProductDetails{
			ProductId:   lFileData.ProductId,
			ProductName: lFileData.ProductName,
			Category:    lFileData.Category,
			UnitPrice:   lFileData.UnitPrice,
			CreatedDate: lFileData.CreatedDate,
			CreatedBy:   lFileData.CreatedBy,
		}

		lOrderRec := models.OrderDetails{
			OrderId:       lFileData.OrderId,
			ProductId:     lFileData.ProductId,
			CustomerId:    lFileData.CustomerId,
			QuantitySold:  lFileData.QuantitySold,
			UnitPrice:     lFileData.UnitPrice,
			Discount:      lFileData.Discount,
			ShippingCost:  lFileData.ShippingCost,
			DateOfSale:    lFileData.DateOfSale,
			Region:        lFileData.Region,
			PaymentMethod: lFileData.PaymentMethod,
			CreatedDate:   lFileData.CreatedDate,
			CreatedBy:     lFileData.CreatedBy,
		}

		lCustomerArr = append(lCustomerArr, lCustomerRec)
		lProductArr = append(lProductArr, lProductRec)
		lOrderArr = append(lOrderArr, lOrderRec)

		// log.Println("lFileData :", lFileData)

		lFileDataArr = append(lFileDataArr, lFileData)
	}

	lResult := lTx.Table("sales_file_data").CreateInBatches(&lFileDataArr, 1000)
	if lResult.Error != nil {
		log.Println("Error IFR 001 :", lResult.Error.Error())
		return lOrderArr, lProductArr, lCustomerArr, lResult.Error
	}

	lErr = lTx.Commit().Error
	if lErr != nil {
		log.Println("Error IFR 002 :", lErr.Error())
		return lOrderArr, lProductArr, lCustomerArr, lErr
	}

	log.Println("InsertFileRecords (-)")
	return lOrderArr, lProductArr, lCustomerArr, nil
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", " ") // Replace non-breaking space
	s = strings.TrimSpace(s)                 // Optional: remove leading/trailing space
	return strings.ToValidUTF8(s, "")        // Remove/replace invalid UTF-8 sequences
}
