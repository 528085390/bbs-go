package db

import (
	"log"
	"temp/common/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDB 获取全局数据库实例（单例模式）
func GetDB() *gorm.DB {

	dsn := "host=localhost user=postgres password=123456 dbname=bbs-go port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	log.Printf("Connecting to database with DSN: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" failed to connect database: %v", err)
	}

	log.Println("✅ Database connected successfully")

	createTable(db)
	log.Println("✅ Database migration completed")

	return db
}

func createTable(db *gorm.DB) {
	var err error
	err = db.AutoMigrate(models.User{})
	err = db.AutoMigrate(models.Section{})
	err = db.AutoMigrate(models.Post{})
	if err != nil {
		log.Fatalf("failed to create tables: " + err.Error())
	}
}
