package models

import (
	"github.com/jinzhu/gorm"
	"github.com/thutasann/book-store/pkg/config"
)

var db *gorm.DB

// Book Model
type Book struct {
	gorm.Model
	Name        string `gorm:"column:name" json:"name"`
	Author      string `string:"author"`
	Publication string `json:"publication"`
}

// Book Model Init
func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

// Create Book
func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

// Get All Books
func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

// Get Book By Id
func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID=?", Id).Find(&getBook)
	return &getBook, db
}

// Delete Book By Id
func DeleteBookById(Id int64) Book {
	var book Book
	db.Where("ID=?", Id).Delete(&book)
	return book
}
