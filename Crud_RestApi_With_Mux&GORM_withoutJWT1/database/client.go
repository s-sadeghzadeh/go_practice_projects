package database

import (
	"log"
	"prj_crud/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var err error

func Connect(connectionString string) {
	Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	//Instance.Set("gorm:auto_preload", true)
	log.Println("Connected to Database...")
}

func Migrate() {
	//Instance.AutoMigrate(&entities.Contact{}, &entities.Company{})
	err = Instance.AutoMigrate(&entities.Contact{}, &entities.Company{})
	if err != nil {
		log.Println("migration failure")
		return
	}
	log.Println("Database Migration Completed...")
}
