package model

type Genre struct {
	ID   int64
	Name string
}

func NewGenre(id int64, name string) *Genre {
	return &Genre{
		ID:   id,
		Name: name,
	}
}
