package filter

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Author struct {
	Name string `json:"name"`
}

type Book struct {
	Name            string `json:"name"`
	Edition         string `json:"edition"`
	PublicationYear uint   `json:"publication_year"`
	AuthorName      string `json:"author_name"`
}

func (a *Author) Filter(db *gorm.DB) *gorm.DB {
	if a.Name != "" {
		db = db.Where("name like ?", "%"+a.Name+"%")
	}
	return db
}

func (b *Book) Filter(db *gorm.DB) *gorm.DB {
	db = db.Preload(clause.Associations)

	if b.Name != "" {
		db = db.Where("name like ?", "%"+b.Name+"%")
	}
	if b.Edition != "" {
		db = db.Where("edition like ?", "%"+b.Edition+"%")
	}
	if b.PublicationYear != 0 {
		db = db.Where("publication_year = ?", b.PublicationYear)
	}
	if b.AuthorName != "" {
		db = db.Joins("inner join book_author ab ON books.id = ab.book_id").
			Joins("inner join authors a ON a.id = ab.author_id").
			Where("a.name like ?", "%"+b.AuthorName+"%")
	}
	return db
}
