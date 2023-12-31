package author

import "gorm.io/gorm"

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
