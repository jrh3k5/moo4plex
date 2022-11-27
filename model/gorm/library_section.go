package gorm

type LibrarySection struct {
	ID   int64  `gorm:"id"`
	Name string `gorm:"name"`
}
