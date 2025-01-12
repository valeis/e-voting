package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"rest-api-go/internal/models"
)

type ConnectionPostgres struct {
	ConnectionString string
}

func (con ConnectionPostgres) Connect() any {
	db, err := gorm.Open(postgres.Open(con.ConnectionString), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Connected to Postgres")
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil
	}
	return db
}
