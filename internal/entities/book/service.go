package book

import (
	"BookApp/internal/entities/author"
	"BookApp/internal/httpResponse"
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

func (s *Service) Create(request Request) (uint, []httpResponse.Cause) {
	var authors []*author.Author
	var causes []httpResponse.Cause
	var result *author.Author

	for _, value := range request.AuthorsId {
		err := s.Db.First(&result, value).Error
		if err != nil {
			cause := httpResponse.Cause{
				Field:   "authors_id",
				Message: fmt.Sprintf("record with id %v not found", value),
			}
			causes = append(causes, cause)
		}
		if result != nil {
			authors = append(authors, result)
		}
		result = nil
	}

	if len(causes) > 0 {
		return 0, causes
	}

	book, err := NewBook(request.Name, request.Edition, request.PublicationYear, authors)
	if err != nil {
		return 0, err
	}

	dbResponseError := s.Db.Create(&book).Error
	if dbResponseError != nil {
		cause := httpResponse.Cause{
			Message: dbResponseError.Error(),
		}

		causes = append(causes, cause)
		return 0, causes
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
