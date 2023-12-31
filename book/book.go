package book

import (
	"BookApp/author"
	"BookApp/pkg/filter"
	"BookApp/pkg/pagination"
	"errors"
	"gorm.io/gorm"
	"time"
)

var (
	InvalidYearInformedErr = errors.New("invalid year informed")
	AuthorsNilErr          = errors.New("authors cannot be null")
)

type Book struct {
	gorm.Model
	Name            string           `json:"name"`
	Edition         string           `json:"edition"`
	PublicationYear uint             `json:"publication_year"`
	Authors         []*author.Author `json:"authors" gorm:"many2many:book_author"`
}

type Request struct {
	Name            string `json:"name"`
	Edition         string `json:"edition"`
	PublicationYear uint   `json:"publication_year"`
	AuthorsId       []uint `json:"authors_id"`
}

type UseCase interface {
	GetById(id uint) (*Book, error)
	Create(request Request) (uint, error)
	Update(id uint, request Request) (uint, error)
	Delete(id uint) error
	List(pagination pagination.Pagination, filter filter.Book) (*pagination.Pagination, error)
}

func NewBook(name, edition string, publicationYear uint, authors []*author.Author) (*Book, error) {
	err := validateBook(publicationYear, authors)
	if err != nil {
		return nil, err
	}

	return &Book{
		Name:            name,
		Edition:         edition,
		PublicationYear: publicationYear,
		Authors:         authors,
	}, nil
}

func validateBook(publicationYear uint, authors []*author.Author) error {
	var errs []error
	now := time.Now()
	if publicationYear > uint(now.Year()) {
		errs = append(errs, InvalidYearInformedErr)
	}
	if authors == nil {
		errs = append(errs, AuthorsNilErr)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (b *Book) UpdateDiffFields(request Request) {
	if request.Name != "" && b.Name != request.Name {
		b.Name = request.Name
	}

	if request.Edition != "" && b.Edition != request.Edition {
		b.Edition = request.Edition
	}

	if request.PublicationYear != 0 && b.PublicationYear != request.PublicationYear {
		b.PublicationYear = request.PublicationYear
	}
}
