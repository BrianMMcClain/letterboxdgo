package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/brianmmcclain/letterboxdgo"
)

func main() {
	slug := flag.String("slug", "", "Film slug (from URL)")
	flag.Parse()

	m := letterboxdgo.GetFilm(*slug)

	fmt.Printf("Title: %s\n", m.Title)
	fmt.Printf("Average Rating: %v\n", m.AvgRating)
	fmt.Printf("Genres: %s\n", strings.Join(m.Genres, ", "))
	fmt.Printf("Rating Count: %v\n", m.Ratings)
	fmt.Printf("Review Count: %v\n", m.Reviews)
	fmt.Printf("TMDB ID: %v\n", m.TMDb)
}
