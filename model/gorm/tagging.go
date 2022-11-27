package gorm

import "time"

type Tagging struct {
	ID             int64     `gorm:"id"`
	MetadataItemID int64     `gorm:"metadata_item_id"`
	TagID          int64     `gorm:"tag_id"`
	Index          int64     `gorm:"index"`
	CreatedAt      time.Time `gorm:"created_at"`
}
