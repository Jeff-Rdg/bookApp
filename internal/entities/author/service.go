package author

import (
	"BookApp/pkg/filter"
	"BookApp/pkg/pagination"
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
)

type Service struct {
	Db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{Db: db}
}

// @Summary Get Author by id
// @Description Get Author by id
// @Tags author
// @Accept json
// @Produce json
// @Param id
// @Success 200 {object} handlers.Response
// @Router /author [get]
func (s *Service) GetById(id uint) (*Author, error) {
	author, err := s.findById(id)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *Service) GetByName(name string) (*Author, error) {
	author, err := s.findByName(name)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *Service) List(pag pagination.Pagination, author filter.Author) (*pagination.Pagination, error) {
	var authors []*Author

	s.Db.Scopes(pagination.Paginate(authors, &pag, "", s.Db, author.Filter)).Find(&authors)
	pag.Rows = authors

	return &pag, nil
}

func (s *Service) Create(author *Author) (uint, error) {
	err := s.Db.Create(&author).Error
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Service) Update(id uint, request Request) (uint, error) {
	author, err := s.GetById(id)
	if err != nil {
		return 0, err
	}
	author.UpdateDiffFields(request)
	err = s.Db.Save(&author).Error
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Service) ReadCsv(reader io.Reader, action func(request Request) error) error {
	csvReader := csv.NewReader(reader)
	_, err := csvReader.Read()
	if err != nil {
		log.Println("error to read first line:", err)
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		request := Request{Name: record[0]}

		err = action(request)
		if err != nil {
			err = fmt.Errorf("author exists: %s", request.Name)
			return err
		}
	}
	return nil
}

// Internal methods
func (s *Service) findById(id uint) (*Author, error) {
	var author *Author
	err := s.Db.First(&author, id).Error
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *Service) findByName(name string) (*Author, error) {
	var author *Author
	err := s.Db.Where("name = ?", name).First(&author).Error
	if err != nil {
		return nil, err
	}
	return author, nil
}
