package gorm

type MetadataItem struct {
	ID    int64  `gorm:"id"`
	Title string `gorm:"title"`
}
