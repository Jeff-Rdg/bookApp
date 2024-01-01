package author

import (
	"encoding/csv"
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

func (s *Service) GetById(id uint) (*Author, error) {
	author, err := s.findById(id)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *Service) Create(request Request) (uint, error) {
	author, err := NewAuthor(request)
	if err != nil {
		return 0, err
	}

	err = s.Db.Create(&author).Error
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

// Internal methods
func (s *Service) findById(id uint) (*Author, error) {
	var author *Author
	err := s.Db.First(&author, id).Error
	if err != nil {
		return nil, err
	}
	return author, nil
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
			log.Printf("author exists: %s", request.Name)
			log.Println(err)
			continue
		}
	}
	return nil
}
