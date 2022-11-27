package model

type Genre struct {
	ID             int64
	Name           string
	MediaLibraryID int64
}

func NewGenre(id int64, name string, mediaLibraryID int64) *Genre {
	return &Genre{
		ID:             id,
		Name:           name,
		MediaLibraryID: mediaLibraryID,
	}
}
