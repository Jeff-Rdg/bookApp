package book

import (
	"BookApp/internal/entities/author"
	"BookApp/internal/httpResponse"
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
	Name            string           `json:"name" gorm:"unique"`
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

func NewBook(name, edition string, publicationYear uint, authors []*author.Author) (*Book, []httpResponse.Cause) {
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

func validateBook(publicationYear uint, authors []*author.Author) []httpResponse.Cause {
	var errs []httpResponse.Cause
	now := time.Now()
	if publicationYear > uint(now.Year()) {
		cause := httpResponse.Cause{
			Field:   "publication_year",
			Message: InvalidYearInformedErr.Error(),
		}
		errs = append(errs, cause)
	}
	if authors == nil {
		cause := httpResponse.Cause{
			Field:   "author_ids",
			Message: AuthorsNilErr.Error(),
		}
		errs = append(errs, cause)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (b *Book) UpdateDiffFields(request Request, db *gorm.DB) error {
	now := time.Now()
	if request.Name != "" && b.Name != request.Name {
		b.Name = request.Name
	}

	if request.Edition != "" && b.Edition != request.Edition {
		b.Edition = request.Edition
	}

	if request.PublicationYear != 0 && b.PublicationYear != request.PublicationYear {
		if request.PublicationYear <= uint(now.Year()) {
			b.PublicationYear = request.PublicationYear
		}
	}

	if len(request.AuthorsId) > 0 {
		var findedAuthors []*author.Author
		db.Model(&b).Association("Authors").Clear()
		service := author.Service{Db: db}
		for _, id := range request.AuthorsId {
			findedAuthor, err := service.GetById(id)
			if err != nil {
				return err
			}
			findedAuthors = append(findedAuthors, findedAuthor)

		}
		db.Model(&b).Association("Authors").Append(findedAuthors)
	}
	return nil
}
