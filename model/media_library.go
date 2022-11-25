package model

// MediaLibrary describes a library of media as it exists within Plex
type MediaLibrary struct {
	Name string // the name of the library as it appears to the user
}

// NewMediaLibrary creates a new instance of MediaLibrary
func NewMediaLibrary(name string) *MediaLibrary {
	return &MediaLibrary{
		Name: name,
	}
}
