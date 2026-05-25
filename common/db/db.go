package db

import (
	"fmt"
	"log"
	"temp/common/db/config"
	"temp/common/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConfig config.DatabaseConfig

func init() {
	dbConfig.SetDefaults()
	dbConfig.LoadFromEnv()
}

func GetDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Dbname, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Database connected successfully")

	return db
}

func CreateTable(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Favorite{},
		&models.Like{},
		&models.Section{},
		&models.Follow{},
		&models.Comment{},
		&models.File{},
	)
	if err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}
}

func Init() {
	CreateTable(GetDB())
}
