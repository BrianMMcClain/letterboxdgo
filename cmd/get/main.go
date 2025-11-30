package main

import (
	"encoding/csv"
	"flag"
	"log/slog"
	"os"
	"strconv"

	"github.com/brianmmcclain/letterboxdgo"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Parse flags
	username := flag.String("username", "", "Letterboxd username")
	minYear := flag.Int("year", 0, "Pull all diary entries from YEAR to current day")
	outFilename := flag.String("out", "stats.csv", "File to write stats to")
	flag.Parse()

	d := letterboxdgo.GetDiary(*username)

	outFile, err := os.Create(*outFilename)
	if err != nil {
		slog.Error("Error creating stats.csv", "error", err)
		os.Exit(-1)
	}
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers := []string{"Title", "ReleaseYear", "Watched", "Rating", "Liked", "Rewatch", "Slug", "TMDB"}
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
		row := []string{
			v.Title,
			strconv.Itoa(v.ReleaseYear),
			v.Date.String(),
			strconv.Itoa(v.Rating),
			strconv.FormatBool(v.Liked),
			strconv.FormatBool(v.Rewatch),
			v.Slug,
			m.TMDb,
		}
		err = writer.Write(row)
		if err != nil {
			slog.Error("Error writing row", "slug", v.Slug, "error", err)
			os.Exit(-1)
		}
	}
}
