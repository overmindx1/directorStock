package models

import (
	"time"
)

type DirectorSelection struct {
	ID        uint
	StockId   string
	StockName string
	AvgCost   string
	Shares    string
	OnList    int
	StockItem *StockItems `gorm:"foreignKey:StockId;references:StockId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (DirectorSelection) TableName() string {
	return "director_selection"
}

type Tabler interface {
	TableName() string
}
