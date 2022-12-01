package gorm

// MetadataType enumerates the known and supported types of metadata
type MetadataType int

const (
	Movie                   MetadataType = 1
	TelevisionSeries        MetadataType = 2
	TelevisionSeriesSeason  MetadataType = 3
	TelevisionSeriesEpisode MetadataType = 4
	AudioArtist             MetadataType = 8
	AudioAlbum              MetadataType = 9
	AudioTrack              MetadataType = 10
)
