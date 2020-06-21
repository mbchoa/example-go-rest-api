package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Book struct representing our Book DB model
type Book struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Author    string    `gorm:"size:100;" json:"author"`
	Title     string    `gorm:"size:100;" json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate verifies that the required fields are present
func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("book: missing required title")
	}
	if b.Author == "" {
		return errors.New("book: missing required author")
	}
	return nil
}

// SaveBook saves the model instance data to the database
func (b *Book) SaveBook(db *gorm.DB) (*Book, error) {
	err := db.Debug().Create(&b).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetAllBooks returns a reference to array of first 100 books in the database
func (b *Book) GetAllBooks(db *gorm.DB) (*[]Book, error) {
	books := []Book{}
	err := db.Debug().Limit(100).Order("id").Find(&books).Error
	if err != nil {
		return nil, err
	}
	return &books, nil
}

// GetBookByID returns a reference to the book given the book ID
func (b *Book) GetBookByID(db *gorm.DB, bid uint64) (*Book, error) {
	err := db.Debug().First(&b, bid).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateBook updates the title and author fields in the database and returns a reference to the updated book
func (b *Book) UpdateBook(db *gorm.DB) (*Book, error) {
	err := db.Debug().Model(b).Updates(Book{Title: b.Title, Author: b.Author}).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteBook removes the selected book from the database
func (b *Book) DeleteBook(db *gorm.DB) (uint64, error) {
	err := db.Debug().Delete(b).Error
	if err != nil {
		return 0, err
	}
	return b.ID, nil
}
