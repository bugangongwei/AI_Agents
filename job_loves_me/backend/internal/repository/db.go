package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/bugangongwei/job_loves_me/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate the schemas
	err = DB.AutoMigrate(&model.User{}, &model.Resume{}, &model.JobDescription{}, &model.Evaluation{})
	if err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
}
