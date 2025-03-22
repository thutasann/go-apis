package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

// Database Connection Func
func Connect() {
	d, err := gorm.Open("mysql", "root:thutasann2002tts/go_bookstore?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d
}

// Get DB
func GetDB() *gorm.DB {
	return db
}
