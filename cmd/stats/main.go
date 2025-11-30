package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/brianmmcclain/letterboxdgo"
)

func main() {

	//slog.SetLogLoggerLevel(slog.LevelDebug)

	inFilename := flag.String("in", "stats.csv", "Source stats file")
	flag.Parse()

	slog.Debug("Reading CSV file", "file", *inFilename)
	inFile, err := os.Open(*inFilename)
	if err != nil {
		slog.Error("Could not open stats file", "file", *inFilename, "error", err)
		os.Exit(-1)
	}
	defer inFile.Close()

	diary := []*letterboxdgo.DiaryEntry{}

	reader := csv.NewReader(inFile)
	// Skip headers
	_, err = reader.Read()
	if err != nil {
		slog.Error("Error reading CSV headers", "error", err)
		os.Exit(-1)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Error reading CSV row", "error", err)
			os.Exit(-1)
		}

		d := &letterboxdgo.DiaryEntry{}
		d.Title = record[0]
		d.ReleaseYear, _ = strconv.Atoi(record[1])
		d.Date, _ = time.Parse("2006-01-02", record[2])
		d.Rating, _ = strconv.Atoi(record[3])
		d.Liked, _ = strconv.ParseBool(record[4])
		d.Rewatch, _ = strconv.ParseBool(record[5])
		d.Slug = record[6]

		diary = append(diary, d)
	}

	count := 0
	liked := 0
	rewatches := 0
	thisYearCount := 0
	watches := make(map[string]int)

	for _, entry := range diary {
		if entry.Date.Year() == 2025 {
			count++
		}
		if entry.Liked {
			liked++
		}
		if entry.Rewatch {
			rewatches++
		}
		if entry.ReleaseYear == 2025 {
			thisYearCount++
		}

		watches[entry.Slug]++
	}

	fmt.Printf("Total entries in 2025: %d\n", count)
	fmt.Printf("Unique films: %d\n", len(slices.Collect(maps.Keys(watches))))
	fmt.Printf("Films liked: %d\n", liked)
	fmt.Printf("Films rewatched: %d\n", rewatches)
	fmt.Printf("2025 films watched: %d\n", thisYearCount)
}
