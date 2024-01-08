package pagination

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"regexp"
	"strconv"
)

var (
	InvalidLimitErr = errors.New("limit must contain only numbers")
	InvalidPageErr  = errors.New("page must contain only numbers")
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func NewPagination(limit, page, sort string) (*Pagination, error) {
	err := validatePagination(limit, page)
	if err != nil {
		return nil, err
	}
	newLimit, newPage := handlePagination(limit, page)

	return &Pagination{
		Limit: newLimit,
		Page:  newPage,
		Sort:  sort,
	}, nil
}

func validatePagination(limit, page string) error {
	var errs []error
	onlyNumbers := regexp.MustCompile("^[0-9]*$")
	if !onlyNumbers.MatchString(limit) {
		errs = append(errs, InvalidLimitErr)
	}
	if !onlyNumbers.MatchString(page) {
		errs = append(errs, InvalidPageErr)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func handlePagination(limit, page string) (int, int) {
	newLimit, err := strconv.Atoi(limit)
	if err != nil {
		newLimit = 0
	}

	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 0
	}

	return newLimit, newPage

}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func Paginate(value interface{}, pag *Pagination, db *gorm.DB, filterFunc func(db *gorm.DB) *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64

	db.Model(value).Scopes(filterFunc).Count(&totalRows)

	pag.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pag.GetLimit())))
	pag.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(filterFunc).Offset(pag.GetOffset()).Limit(pag.GetLimit()).Order(pag.GetSort())
	}
}
