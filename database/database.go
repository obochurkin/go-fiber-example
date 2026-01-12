package database

import (
	"log"

	"github.com/obochurkin/go-fiber-example/config"
	"github.com/obochurkin/go-fiber-example/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DB *gorm.DB
}

var Instance DBInstance

func Connect() {
	db, err := gorm.Open(postgres.Open(config.GetEnvVariable("PG_DNS")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln("Failed to connect to database")
	}

	log.Println("Database connection established")

	Instance = DBInstance{DB: db}

	SyncSchema()
}

func SyncSchema() {
	log.Println("Syncing database schema...")
	err := Instance.DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatalln("Failed to sync database schema")
	}
	log.Println("Database schema synced successfully")
}
