package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbUri = "root:change-me@tcp(localhost:3307)/gofibersampledb?charset=utf8&parseTime=True&loc=Local"

func connectDb() error {
	var err error
	db, err = gorm.Open(mysql.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
