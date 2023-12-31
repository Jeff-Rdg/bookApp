package author

import (
	"errors"
	"gorm.io/gorm"
	"regexp"
)

var (
	InvalidNameLenErr = errors.New("invalid name length")
	InvalidNameErr    = errors.New("invalid name")
)

type Author struct {
	gorm.Model
	Name string `json:"name"`
}

type Request struct {
	Name string `json:"name"`
}

type UseCase interface {
	GetById(id uint) (*Author, error)
	Create(request Request) (uint, error)
	Update(id uint, request Request) (uint, error)
}

// constructors
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
	regexName := regexp.MustCompile("^[a-zA-ZÀ-ÖØ-öø-ÿ\\s]+$")

	return regexName.MatchString(name)
}
