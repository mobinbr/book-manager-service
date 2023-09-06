package db

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID                int     `gorm:"primaryKey;autoIncrement"`
	Name              string   `gorm:"varchar(50)"`
	AuthorFirstName   string   `gorm:"varchar(50)"` // author
	AuthorLastName    string   `gorm:"varchar(50)"` // author
	AuthorBirthday    string   `gorm:"varchar(30)"` // author
	AuthorNationality string   `gorm:"varchar(50)"` // author
	Category          string   `gorm:"varchar(50)"`
	Volume            int      `gorm:"int"`
	PublishedAt       string   `gorm:"varchar(30)"`
	Summary           string   `gorm:"varchar(100)"`
	TableOfContents   []string `gorm:"json"`
	Publisher         string   `gorm:"varchar(50)"`
	UserOwner         string   `gorm:"varchar(50)"`
}

// AddNewBook adds a books to the database using GormDB
func (gdb *GormDB) AddNewBook(b *Book) error {
	return gdb.db.Create(b).Error
}

// GetBooksByUsername returns the books of the given username
func (gdb *GormDB) GetBooksByUsername(username string) (*[]Book, error) {
	var books []Book
	err := gdb.db.Where(&Book{UserOwner: username}).Find(&books).Error
	if err != nil {
		return nil, err
	}

	return &books, nil
}

func (gdb *GormDB) GetBookByID(id int) (*Book, error) {
	var book Book
	err := gdb.db.Where(&Book{ID: id}).First(&book).Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (gdb *GormDB) DeleteBookByID(id int) error {
	err := gdb.db.Delete(&Book{ID: id}).Error
	if err != nil {
		return err
	}
	
	return nil
}

func (gdb *GormDB) UpdateBook(book *Book) error {
	err := gdb.db.Save(book).Error
	if err != nil {
		return err
	}
	
	return nil
}