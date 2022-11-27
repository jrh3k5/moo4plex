package model

// MediaLibrary describes a library of media as it exists within Plex
type MediaLibrary struct {
	ID   int64  // the ID of the library
	Name string // the name of the library as it appears to the user
}

// NewMediaLibrary creates a new instance of MediaLibrary
func NewMediaLibrary(id int64, name string) *MediaLibrary {
	return &MediaLibrary{
		ID:   id,
		Name: name,
	}
}
