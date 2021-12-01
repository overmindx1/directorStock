package modelSocket

import (
	"time"

	"github.com/overmind/go_test/models"
)

type StockItemsNoMem struct {
	ID                uint
	StockId           string
	StockName         string
	StockType         string
	OpenPrice         string
	ClosePrice        string
	UpDown            string
	UpDownCount       string
	Date              string
	DirectorSelection *models.DirectorSelection `gorm:"foreignKey:StockId;references:StockId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (StockItemsNoMem) TableName() string {
	return "stock_items"
}

type Tabler interface {
	TableName() string
}
