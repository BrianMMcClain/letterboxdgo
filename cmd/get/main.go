package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/brianmmcclain/letterboxdgo"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	username := os.Args[1]
	d := letterboxdgo.GetDiary(username)

	for _, v := range d {
		//fmt.Printf("%+v\n", v)
		fmt.Printf("%s,%s,%d,%t,%t,%d\n", v.Title, v.Date.Format("2006-01-02"), v.Rating, v.Liked, v.Rewatch, v.ReleaseYear)
	}
}
