package gorm

type Tag struct {
	ID               int64  `gorm:"id"`                 // the ID of the tag
	Tag              string `gorm:"tag"`                // the name of the tag
	TagType          int64  `gorm:"tag_type"`           // the type of the tag
	LibrarySectionID int64  `gorm:"library_section_id"` // the ID of the library section to which the tag belongs
	UserThumbURL     string `gorm:"user_thumb_url"`     // the URL of the thumb art image for a particular tag, if available and relevant
}
