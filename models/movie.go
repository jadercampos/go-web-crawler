package models

type Movie struct {
	Title string
	Year  string
}
type Star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []Movie
}
