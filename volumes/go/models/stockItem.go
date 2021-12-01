package models

import (
	"time"
)

type StockItems struct {
	ID                uint
	StockId           string
	StockName         string
	StockType         string
	OpenPrice         string
	ClosePrice        string
	UpDown            string
	UpDownCount       string
	Date              string
	DirectorSelection *DirectorSelection `gorm:"foreignKey:StockId;references:StockId"`
	MemberData        *MemberData        `gorm:"embedded"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
