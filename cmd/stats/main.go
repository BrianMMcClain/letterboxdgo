package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/brianmmcclain/letterboxdgo"
)

type Movie struct {
	Entry    *letterboxdgo.DiaryEntry
	Genres   []string
	Language string
	Runtime  int
}

func main() {

	//slog.SetLogLoggerLevel(slog.LevelDebug)

	inFilename := flag.String("in", "stats.csv", "Source stats file")
	generateGraphs := flag.Bool("graphs", false, "Generate graphs for blog")
	flag.Parse()

	slog.Debug("Reading CSV file", "file", *inFilename)
	inFile, err := os.Open(*inFilename)
	if err != nil {
		slog.Error("Could not open stats file", "file", *inFilename, "error", err)
		os.Exit(-1)
	}
	defer inFile.Close()

	movies := []*Movie{}
	ratings := []int{}
	totalRating := 0

	reader := csv.NewReader(inFile)
	// Skip headers
	_, err = reader.Read()
	if err != nil {
		slog.Error("Error reading CSV headers", "error", err)
		os.Exit(-1)
	}

	var first *Movie = nil
	var last *Movie = nil
	var oldest *Movie = nil

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Error reading CSV row", "error", err)
			os.Exit(-1)
		}

		m := new(Movie)
		entry := &letterboxdgo.DiaryEntry{}
		entry.Title = record[0]
		entry.ReleaseYear, _ = strconv.Atoi(record[1])
		entry.Date, _ = time.Parse("2006-01-02", record[2])
		entry.Rating, _ = strconv.Atoi(record[3])
		entry.Liked, _ = strconv.ParseBool(record[4])
		entry.Rewatch, _ = strconv.ParseBool(record[5])
		entry.Slug = record[6]
		m.Entry = entry

		m.Genres = strings.Split(record[8], "|")
		m.Language = record[9]
		m.Runtime, _ = strconv.Atoi(record[10])

		movies = append(movies, m)
		ratings = append(ratings, entry.Rating)
		totalRating += entry.Rating

		if first == nil || first.Entry.Date.After(entry.Date) || first.Entry.Date.Equal(entry.Date) {
			first = m
		}
		if last == nil || last.Entry.Date.Before(entry.Date) {
			last = m
		}
		if oldest == nil || oldest.Entry.ReleaseYear > entry.ReleaseYear {
			oldest = m
		}
	}

	count := 0
	liked := 0
	rewatches := 0
	thisYearCount := 0
	watches := make(map[string]int)
	ratingDistribution := map[int]int{
		0:  0,
		1:  0,
		2:  0,
		3:  0,
		4:  0,
		5:  0,
		6:  0,
		7:  0,
		8:  0,
		9:  0,
		10: 0,
	}

	runtime := 0

	for _, movie := range movies {
		if movie.Entry.Date.Year() == 2025 {
			count++
		}
		if movie.Entry.Liked {
			liked++
		}
		if movie.Entry.Rewatch {
			rewatches++
		}
		if movie.Entry.ReleaseYear == 2025 {
			thisYearCount++
		}

		ratingDistribution[movie.Entry.Rating]++

		watches[movie.Entry.Slug]++
		runtime += movie.Runtime
	}

	runtimeHours := int(runtime / 60)
	runtimeMinutes := runtime % 60

	averageRating := float64(totalRating) / float64(count)

	// Calculate median rating
	medianRating := 0
	slices.Sort(ratings)
	if len(ratings)%2 == 0 {
		// Even number of elements, find the average of the two middle numbers
		lIndex := int(math.Floor(float64(len(ratings)) / 2.0))
		rIndex := int(math.Ceil(float64(len(ratings)) / 2.0))
		medianRating = (ratings[lIndex] + ratings[rIndex]) / 2.0
	} else {
		medianRating = ratings[len(ratings)/2]
	}

	fmt.Printf("Total entries in 2025: %d\n", count)
	fmt.Printf("Unique films: %d\n", len(slices.Collect(maps.Keys(watches))))
	fmt.Printf("Films liked: %d\n", liked)
	fmt.Printf("Films rewatched: %d\n", rewatches)
	fmt.Printf("2025 films watched: %d\n", thisYearCount)
	fmt.Printf("Total runtime: %dh %dm\n", runtimeHours, runtimeMinutes)
	fmt.Println()
	fmt.Printf("Avg rating: %.1f\n", averageRating)
	fmt.Printf("Median rating: %d\n", medianRating)
	fmt.Println()
	fmt.Printf("First movie watched: %s\n", first.Entry.Title)
	fmt.Printf("Last movie watched: %s\n", last.Entry.Title)
	fmt.Printf("Oldest movie watched: %s (%d)\n", oldest.Entry.Title, oldest.Entry.ReleaseYear)

	fmt.Println("\n\n")
	if *generateGraphs {
		// Collect values
		ratingDistributionValues := []string{}
		for i := 0; i <= 10; i++ {
			ratingDistributionValues = append(ratingDistributionValues, strconv.Itoa(ratingDistribution[i]))
		}

		// Print graphs
		// Rating distribution
		ratingDistributionGraph := fmt.Sprintf(`{{< chart id="allRatings" title="Rating distribution of all watched films" width="100" >}}

{
 type: 'bar',
 data: {
    labels: ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '10'],
    datasets: [{
        label: '# of Films',
        data: [%s],
        borderWidth: 1
    }]
 },
 options: {
    scales: {
        y: {
            beginAtZero: true
        }
    }
 }
};

{{</ chart >}}`, strings.Join(ratingDistributionValues, ","))

		fmt.Println(ratingDistributionGraph)
		fmt.Println("\n\n")
		// Table
		fmt.Println("|Rating|Count|")
		fmt.Println("|------|-----|")
		for i := 0; i <= 10; i++ {
			fmt.Printf("|%d|%d|\n", i, ratingDistribution[i])
		}
	}
}
