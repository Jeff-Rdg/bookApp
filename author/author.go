package author

import (
	"BookApp/pkg/filter"
	"BookApp/pkg/pagination"
	"errors"
	"gorm.io/gorm"
	"io"
	"regexp"
)

var (
	InvalidNameLenErr = errors.New("invalid name length")
	InvalidNameErr    = errors.New("invalid name")
)

type Author struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}

type Request struct {
	Name string `json:"name"`
}

type UseCase interface {
	GetById(id uint) (*Author, error)
	GetByName(name string) ([]*Author, error)
	Create(request Request) (uint, error)
	Update(id uint, request Request) (uint, error)
	ReadCsv(reader io.Reader, action func(request Request) error) error
	List(pagination pagination.Pagination, filter filter.Author) (*pagination.Pagination, error)
}

// NewAuthor receiver Request with params, and create a new User
func NewAuthor(request Request) (*Author, error) {
	err := validateUser(request.Name)
	if err != nil {
		return nil, err
	}

	return &Author{
		Name: request.Name,
	}, nil
}

// Enterprise Rules
func validateUser(name string) error {
	var errs []error
	if len(name) < 3 {
		errs = append(errs, InvalidNameLenErr)
	}
	if !isValidName(name) {
		errs = append(errs, InvalidNameErr)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func isValidName(name string) bool {
	regexName := regexp.MustCompile("^[a-zA-ZÀ-ÖØ-öø-ÿ\\s.]+$")

	return regexName.MatchString(name)
}

func (a *Author) UpdateDiffFields(request Request) {
	if request.Name != "" {
		if a.Name != request.Name {
			a.Name = request.Name
		}
	}
}
