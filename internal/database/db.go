package database

import (
	"fmt"

	"github.com/simple-crud-go/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", configs.GetDBUSER(), configs.GetDBPASS(), configs.GetDBHOST(), configs.GetDBPORT(), configs.GetDBNAME()),
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db
	return db
}

func GetDBGorm() *gorm.DB {
	return DB
}
