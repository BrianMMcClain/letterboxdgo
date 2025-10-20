package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brianmmcclain/letterboxdgo"
)

func main() {

	file, err := os.Open("./data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	diary := []*letterboxdgo.DiaryEntry{}

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		d := &letterboxdgo.DiaryEntry{}
		d.Title = record[0]
		d.Date, _ = time.Parse("2006-01-02", record[1])
		d.Rating, _ = strconv.Atoi(record[2])
		d.Liked, _ = strconv.ParseBool(record[3])
		d.Rewatch, _ = strconv.ParseBool(record[4])
		d.ReleaseYear, _ = strconv.Atoi(record[5])

		diary = append(diary, d)
	}

	count := 0
	for _, entry := range diary {
		if entry.Date.Year() == 2025 {
			log.Printf("%+v\n", entry)
			count++
		}
	}

	log.Printf("Total entries in 2025: %d\n", count)
}
