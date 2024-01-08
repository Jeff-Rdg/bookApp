package filter

import "gorm.io/gorm"

type Author struct {
	Name string `json:"name"`
}

func (a *Author) Filter(db *gorm.DB) *gorm.DB {
	if a.Name != "" {
		db = db.Where("name like ?", "%"+a.Name+"%")
	}
	return db
}
