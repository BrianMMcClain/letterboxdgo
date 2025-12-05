package main

import (
	"encoding/csv"
	"flag"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/brianmmcclain/letterboxdgo"
	"github.com/brianmmcclain/tmdbgo"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Parse flags
	username := flag.String("username", "", "Letterboxd username")
	minYear := flag.Int("year", 0, "Pull all diary entries from YEAR to current day")
	outFilename := flag.String("out", "stats.csv", "File to write stats to")
	flag.Parse()

	tmdbKey := os.Getenv("TMDB_KEY")
	includeTMDB := tmdbKey != ""

	tmdb := new(tmdbgo.TMDB)
	if includeTMDB {
		tmdb = tmdbgo.NewTMDB(tmdbKey)
	}

	d := letterboxdgo.GetDiary(*username)

	outFile, err := os.Create(*outFilename)
	if err != nil {
		slog.Error("Error creating stats.csv", "error", err)
		os.Exit(-1)
	}
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers := []string{"Title", "ReleaseYear", "Watched", "Rating", "Liked", "Rewatch", "Slug", "TMDB", "Generes", "Language", "Runtime", "AvgRating", "Ratings", "Reviews"}
	err = writer.Write(headers)
	if err != nil {
		slog.Error("Error writing headers", "error", err)
		os.Exit(-1)
	}

	for _, v := range d {

		if v.Date.Year() < *minYear {
			continue
		}

		m := letterboxdgo.GetFilm(v.Slug)

		genres := []string{}
		language := ""
		runtime := 0
		if includeTMDB {
			tmdbMovie := tmdb.GetMovie(m.TMDb)
			for _, g := range tmdbMovie.Genres {
				genres = append(genres, g.Name)
			}
			language = tmdbMovie.Language
			runtime = tmdbMovie.Runtime
		}

		row := []string{
			v.Title,
			strconv.Itoa(v.ReleaseYear),
			v.Date.Format("2006-01-02"),
			strconv.Itoa(v.Rating),
			strconv.FormatBool(v.Liked),
			strconv.FormatBool(v.Rewatch),
			v.Slug,
			m.TMDb,
			strings.Join(genres, "|"),
			language,
			strconv.Itoa(runtime),
			strconv.FormatFloat(float64(m.AvgRating), 'f', 2, 64),
			strconv.Itoa(int(m.Ratings)),
			strconv.Itoa(int(m.Reviews)),
		}
		err = writer.Write(row)
		if err != nil {
			slog.Error("Error writing row", "slug", v.Slug, "error", err)
			os.Exit(-1)
		}
	}
}
