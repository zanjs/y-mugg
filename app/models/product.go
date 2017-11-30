package models

import (
	"time"

	"github.com/zanjs/y-mugg/db"
)

type (
	// Product is
	Product struct {
		BaseModel
		Title        string `json:"title" gorm:"type:varchar(100)"`
		ExternalCode string `json:"external_code" gorm:"varchar(100);not null;unique"`
		Sort         int    `json:"sort"`
	}
	// QMProduct is
	QMProduct struct {
		OwnerCode     string `json:"ownerCode"`
		ItemCode      string `json:"itemCode"`
		WarehouseCode string `json:"warehouseCode"`
		InventoryType string `json:"inventoryType"`
	}
)

// CreateProduct is
func CreateProduct(m *Product) error {
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

// UpdateProduct is
func (m *Product) UpdateProduct(data *Product) error {
	var err error

	m.UpdatedAt = time.Now()
	m.Title = data.Title
	m.ExternalCode = data.ExternalCode
	m.Sort = data.Sort

	tx := gorm.MysqlConn().Begin()
	if err = tx.Save(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

// DeleteProduct is
func (m Product) DeleteProduct() error {
	var err error
	tx := gorm.MysqlConn().Begin()
	if err = tx.Delete(&m).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}

// GetProductById is
func GetProductById(id uint64) (Product, error) {
	var (
		product Product
		err     error
	)

	tx := gorm.MysqlConn().Begin()
	if err = tx.Last(&product, id).Error; err != nil {
		tx.Rollback()
		return product, err
	}
	tx.Commit()

	return product, err
}

// GetProducts is
func GetProducts() ([]Product, error) {
	var (
		products []Product
		err      error
	)

	tx := gorm.MysqlConn().Begin()
	if err = tx.Order("sort desc").Find(&products).Error; err != nil {
		tx.Rollback()
		return products, err
	}
	tx.Commit()

	return products, err
}
