package gorm

type Tag struct {
	ID      int64  `gorm:"id"`       // the ID of the tag
	Tag     string `gorm:"tag"`      // the name of the tag
	TagType int64  `gorm:"tag_type"` // the type of the tag
}
