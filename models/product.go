package models

import "time"

type Product struct {
	ID         int    `json:"id,string"`
	Title      string `json:"title" gorm:"unique"`
	CreatedBy  int
	User       User `gorm:"foreignKey:CreatedBy" json:"-"`
	CategoryId int  `json:"category_id,string"`
	Price      int  `json:"price"`
	// Meta
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
