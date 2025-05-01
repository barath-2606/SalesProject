package processsalesfile

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm/clause"
	"main.go/db"
	"main.go/models"
)

func ProcessCustomerDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("ProcessCustomerDetails (+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", http.MethodPost)
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSTF-Token,Authorization")

	if !(strings.EqualFold(http.MethodPost, r.Method)) {
		if strings.EqualFold(http.MethodOptions, r.Method) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	var lReq models.CustomerReq
	var lResp models.ReadDataResp
	lResp.Status = "S"

	lBody, lErr := io.ReadAll(r.Body)
	if lErr != nil {
		log.Println("Error PCD 001 :", lErr.Error())
		lResp.Status = "E"
		lResp.ErrMsg = "Error PCD 001 : " + lErr.Error()
	}

	lErr = json.Unmarshal(lBody, &lReq)
	if lErr != nil {
		log.Println("Error PCD 002 :", lErr.Error())
		lResp.Status = "E"
		lResp.ErrMsg = "Error PCD 002 : " + lErr.Error()
	}

	lErr = ProcessRequest(lReq.ReqArr)
	if lErr != nil {
		log.Println("Error PCD 003 :", lErr.Error())
		lResp.Status = "E"
		lResp.ErrMsg = "Error PCD 003 : " + lErr.Error()
	}

	lData, lErr := json.Marshal(lResp)
	if lErr != nil {
		log.Println("Error PCD 004 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("ProcessCustomerDetails (-)")
}

func ProcessRequest(pReqArr interface{}) error {
	log.Println("ProcessRequest (+)")

	lTx := db.GPostgres.Begin()
	defer lTx.Rollback()
	lResult := lTx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(pReqArr, 1000)
	if lResult.Error != nil {
		log.Println("Error PREQ 001:", lResult.Error.Error())
		return lResult.Error
	}

	if lErr := lTx.Commit().Error; lErr != nil {
		log.Println("Error PREQ 002 :", lErr.Error())
		return lErr
	}

	log.Println("ProcessRequest (-)")
	return nil
}

func FetchRequestedData(pReqDet interface{}) {
	log.Println("FetchRequestedData (+)")

	log.Println("FetchRequestedData (-)")
}

func RefreshDailySalesSummary() error {
	const refreshQuery = `REFRESH MATERIALIZED VIEW CONCURRENTLY daily_sales_summary`

	if err := db.GPostgres.Exec(refreshQuery).Error; err != nil {
		log.Println("Error refreshing materialized view:", err)
		return err
	}

	log.Println("Materialized view 'daily_sales_summary' refreshed successfully.")
	return nil
}
