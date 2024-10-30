package models

import "time"

type Books struct {
	ID            uint      `gorm:primaryKey json:"id"`
	Title         string    `gorm:not null json:"title"`
	Description   string    `gorm:not null json:"description"`
	Author        string    `gorm:not null json:"author"`
	Publisher     string    `gorm: not null json:"publisher"`
	CoverImage    string    `gorm: null json:"coverimage"`
	BookUrl       string    `gorm: null json:"bookurl"`
	Rating        int       `gorm: nulljson:"rating"`
	Category      string    `gorm: null json:"category"`
	PublisherName string    `json:"publishername"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"-"`
}
