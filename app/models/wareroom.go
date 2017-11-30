package models

import (
	"time"

	"github.com/zanjs/y-mugg/db"
)

type (
	// Wareroom is
	Wareroom struct {
		BaseModel
		Title     string `json:"title" gorm:"type:varchar(100)"`
		Numbering string `json:"numbering" gorm:"type:varchar(100);not null;unique"`
		Sort      int    `json:"sort"`
	}
)

// CreateWareroom is
func CreateWareroom(m *Wareroom) error {
	var err error

	m.CreatedAt = time.Now()
	tx := gorm.MysqlConn().Begin()
	if err = tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

// UpdateWareroom is
func (m *Wareroom) UpdateWareroom(data *Wareroom) error {
	var err error

	m.UpdatedAt = time.Now()
	m.Title = data.Title
	m.Numbering = data.Numbering
	m.Sort = data.Sort

	tx := gorm.MysqlConn().Begin()
	if err = tx.Save(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

// DeleteWareroom is
func (m Wareroom) DeleteWareroom() error {
	var err error
	tx := gorm.MysqlConn().Begin()
	if err = tx.Delete(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

// GetWareroomById is
func GetWareroomById(id uint64) (Wareroom, error) {
	var (
		wareroom Wareroom
		err      error
	)

	tx := gorm.MysqlConn().Begin()
	if err = tx.Last(&wareroom, id).Error; err != nil {
		tx.Rollback()
		return wareroom, err
	}
	tx.Commit()

	return wareroom, err
}

func GetWarerooms() ([]Wareroom, error) {
	var (
		warerooms []Wareroom
		err       error
	)

	tx := gorm.MysqlConn().Begin()
	if err = tx.Order("sort desc").Find(&warerooms).Error; err != nil {
		tx.Rollback()
		return warerooms, err
	}
	tx.Commit()

	return warerooms, err
}
