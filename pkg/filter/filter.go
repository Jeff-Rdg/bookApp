package filter

import "gorm.io/gorm"

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
	db = db.Preload("Author")
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
		db = db.Where("authors.name like ?", "%"+b.AuthorName+"%")
	}
	return db
}
