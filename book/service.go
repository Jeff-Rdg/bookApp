package book

import (
	"BookApp/author"
	"BookApp/pkg/filter"
	"BookApp/pkg/pagination"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	Db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{Db: db}
}

func (s *Service) GetById(id uint) (*Book, error) {
	book, err := s.findById(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *Service) Create(request Request) (uint, error) {
	var authors []*author.Author
	var result *author.Author
	for _, value := range request.AuthorsId {
		err := s.Db.First(&result, value).Error
		if err != nil {
			return 0, fmt.Errorf("record with id %v not found", value)
		}
		authors = append(authors, result)
		result = nil
	}
	book, err := NewBook(request.Name, request.Edition, request.PublicationYear, authors)
	if err != nil {
		return 0, err
	}

	err = s.Db.Create(&book).Error
	if err != nil {
		return 0, err
	}

	return book.ID, nil
}

func (s *Service) Update(id uint, request Request) (uint, error) {
	book, err := s.GetById(id)
	if err != nil {
		return 0, err
	}

	err = book.UpdateDiffFields(request, s.Db)
	if err != nil {
		return 0, err
	}

	err = s.Db.Save(&book).Error
	if err != nil {
		return 0, err
	}

	return book.ID, nil
}

func (s *Service) Delete(id uint) error {
	err := s.Db.Delete(&Book{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) List(pag pagination.Pagination, book filter.Book) (*pagination.Pagination, error) {
	var books []*Book
	s.Db.Scopes(pagination.Paginate(books, &pag, "books", s.Db, book.Filter)).Find(&books)
	pag.Rows = books

	return &pag, nil
}

// Internal methods
func (s *Service) findById(id uint) (*Book, error) {
	var book *Book
	err := s.Db.Preload(clause.Associations).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}
