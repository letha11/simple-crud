package database

import (
	"fmt"

	"github.com/simple-crud-go/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%v:%v@tcp(%v:%v)/", configs.GetDBUSER(), configs.GetDBPASS(), configs.GetDBHOST(), configs.GetDBPORT()),
	), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	_ = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v;", configs.GetDBNAME()))

	normalDSN := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", configs.GetDBUSER(), configs.GetDBPASS(), configs.GetDBHOST(), configs.GetDBPORT(), configs.GetDBNAME())
	db, err = gorm.Open(
		mysql.Open(normalDSN),
		&gorm.Config{},
	)

	if err != nil {
		panic(err.Error())
	}

	DB = db
	return db
}

func GetDBGorm() *gorm.DB {
	return DB
}
