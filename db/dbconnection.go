package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Type      string `json:"type"`
	URI       string `json:"uri"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Timeout   int    `json:"timeout"`
	Profile   bool   `json:"profile"`
	RunScript bool   `json:"runScript"`
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Version  string         `json:"version"`
	Profile  bool           `json:"profile"`
	Database DatabaseConfig `json:"database"`
}

var GPostgres *gorm.DB
var ConfigData Config

func LoadConfig(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &ConfigData)
	if err != nil {
		return err
	}

	return nil
}

func ConnectDB() error {
	log.Println("ConnectDB (+)")

	// Fetching database details from loaded config
	dbConfig := ConfigData.Database
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.URI, 5432, dbConfig.Username, dbConfig.Password, dbConfig.Name)

	var lErr error
	GPostgres, lErr = gorm.Open(postgres.Open(dbConnStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if lErr != nil {
		log.Println("Error CDB 001 :", lErr.Error())
		return lErr
	}

	log.Println("DbConnection Established")
	log.Println("ConnectDB (-)")
	return nil
}
