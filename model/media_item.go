package model

// MediaItem is a representation of a piece of media in Plex (a movie, a song, an ablum)
type MediaItem struct {
	ID             int64  // the ID of the media item
	Name           string // the name of the media item
	MediaLibraryID int64  // the ID of the media library to which the item belongs
}

// NewMediaItem creates a new instance of MediaItem
func NewMediaItem(id int64, name string, mediaLibraryID int64) *MediaItem {
	return &MediaItem{
		ID:             id,
		Name:           name,
		MediaLibraryID: mediaLibraryID,
	}
}
