package entities

type Location struct {
	ID          int
	Name        string
	Description string
	Geom        string
	Longitude   *string
	Latitude    *string
}
