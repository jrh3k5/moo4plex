package gorm

type MetadataItem struct {
	ID               int64  `gorm:"id"`
	Title            string `gorm:"title"`
	LibrarySectionID int64  `gorm:"library_section_id"`
}
