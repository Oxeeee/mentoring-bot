package db

import (
	"log"
	"os"

	"github.com/Oxeeee/klenov-bot/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	if _, err := os.Stat("data.db"); os.IsNotExist(err) {
		file, err := os.Create("data.db")
		if err != nil {
			log.Fatalf("Error while creating DB: %v", err)
		}
		file.Close()
	}

	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}
	err = DB.AutoMigrate(&domain.User{}, &domain.Message{})
	if err != nil {
		log.Fatalf("Error while migrating tables: %v", err)
	}

	log.Println("Connected to database")
}

func CreateDefaultAdmin() {
	admin := domain.User{
		Username:      "petrushin_leonid",
		IsWhitelisted: true,
		Role:          "admin",
	}

	var exitingAdmin domain.User
	res := DB.Where("username = ?", admin.Username).First(&exitingAdmin)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Fatalf("Error while searching administrator: %v", res.Error)
	}

	if res.RowsAffected == 0 {
		DB.Create(&admin)
		log.Println("Default administrator created")
	} else {
		log.Println("Administrator already exists")
	}
}
