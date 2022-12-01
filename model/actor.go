package model

// Actor is a representation of actor data in Plex
type Actor struct {
	ID           int64  // the ID of the actor
	Name         string // the name of the actor
	ThumbnailURL string // the URL of the actor's thumbnail data
}

// NewActor creates a new instance of Actor
func NewActor(id int64, name string, thumbnailURL string) *Actor {
	return &Actor{
		ID:           id,
		Name:         name,
		ThumbnailURL: thumbnailURL,
	}
}
