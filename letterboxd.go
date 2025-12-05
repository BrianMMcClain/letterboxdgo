package letterboxdgo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DiaryEntry struct {
	Title       string
	Slug        string
	Date        time.Time
	Rating      int
	Rewatch     bool
	Liked       bool
	ReleaseYear int
}

type Film struct {
	Title     string
	TMDb      string
	AvgRating float64
	Ratings   int
	Reviews   int
	Genres    []string
}

type FilmData struct {
	AggregateRating struct {
		AvgRating float64 `json:"ratingValue"`
		Ratings   int     `json:"ratingCount"`
		Reviews   int     `json:"reviewCount"`
	} `json:"aggregateRating"`
	Genres []string `json:"genre"`
}

func GetFilm(slug string) *Film {

	f := new(Film)

	slog.Debug("Sending request to get film", "slug", slug)
	res, err := Get(fmt.Sprintf("https://letterboxd.com/film/%s", slug))
	if err != nil {
		slog.Error("Error getting film", "slug", slug, "error", err)
		os.Exit(1)
	}

	defer res.Close()
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		slog.Error("Error parsing film page", "slug", slug, "error", err)
		os.Exit(1)
	}

	// Get title
	f.Title = doc.Find("div.details span.name").Text() //.Each(func(i int, s *goquery.Selection) {

	// Get TMDb ID
	f.TMDb, _ = doc.Find("body").Attr("data-tmdb-id")

	// Get script JSON
	jData := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(doc.Find("script[type=\"application/ld+json\"]").Text(), "/* <![CDATA[ */", ""), "/* ]]> */", ""))
	var filmData FilmData
	err = json.Unmarshal([]byte(jData), &filmData)
	if err != nil {
		slog.Error("Error parsing film JSON data", "slug", slug, "error", err)
		os.Exit(1)
	}
	f.AvgRating = filmData.AggregateRating.AvgRating
	f.Ratings = filmData.AggregateRating.Ratings
	f.Reviews = filmData.AggregateRating.Reviews
	f.Genres = filmData.Genres

	return f
}

func GetDiary(user string) []*DiaryEntry {
	// Get page count
	pageCount := 0

	slog.Debug("Sending request to get diary page")
	res, err := Get(fmt.Sprintf("https://letterboxd.com/%s/diary/", user))
	if err != nil {
		slog.Error("Error getting diary page", "error", err)
		os.Exit(1)
	}

	defer res.Close()
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		slog.Error("Error parsing diary page", "error", err)
	}

	doc.Find("li.paginate-page").Each(func(i int, s *goquery.Selection) {
		pageCount = i
	})
	pageCount++

	var diary []*DiaryEntry
	for i := 1; i <= pageCount; i++ {
		slog.Debug("Getting diary page", "count", i)
		page := getDiaryPage(user, i)
		diary = append(diary, page...)
	}

	// Sort entire diary by watch date
	sort.Slice(diary, func(i, j int) bool {
		return diary[i].Date.After(diary[j].Date)
	})

	return diary
}

func getDiaryPage(user string, page int) []*DiaryEntry {

	url := fmt.Sprintf("https://letterboxd.com/%s/diary/films/page/%v", user, page)
	slog.Debug("Sending request to get diary page", "page", page)
	res, err := Get(url)
	if err != nil {
		slog.Error("Error getting diary page", "error", err)
		os.Exit(1)
	}

	defer res.Close()
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		slog.Error("Error parsing diary page", "error", err)
	}

	var entries []*DiaryEntry

	doc.Find("tr.diary-entry-row").Each(func(i int, s *goquery.Selection) {
		entry := new(DiaryEntry)

		// Get title and slug
		s.Find("td.col-production").Each(func(i int, r *goquery.Selection) {
			entry.Title = r.Find("h2.name").Find("a").Text()
			entry.Slug, _ = r.Find("div.react-component").Attr("data-item-slug")
		})

		// Get watch date
		s.Find("td.col-daydate a").Each(func(i int, r *goquery.Selection) {
			fullDatePath, _ := r.Attr("href")
			datePath := strings.ReplaceAll(fullDatePath, fmt.Sprintf("/%s/diary/films/for/", user), "")
			entry.Date, _ = time.Parse("2006/01/02/", datePath)
		})

		// Get rating
		s.Find("input.rateit-field").Each(func(i int, r *goquery.Selection) {
			sRating, _ := r.Attr("value")
			entry.Rating, _ = strconv.Atoi(sRating)
		})

		// Get Rewatch
		s.Find("td.col-rewatch").Each(func(i int, r *goquery.Selection) {
			classes, _ := r.Attr("class")
			if strings.Contains(classes, "icon-status-off") {
				entry.Rewatch = false
			} else {
				entry.Rewatch = true
			}
		})

		// Get Liked
		entry.Liked = false
		s.Find("td.col-like").Each(func(i int, r *goquery.Selection) {
			r.Find("span.icon-liked").Each(func(i int, l *goquery.Selection) {
				entry.Liked = true
			})
		})

		// Get release year
		s.Find("td.col-releaseyear").Each(func(i int, r *goquery.Selection) {
			entry.ReleaseYear, _ = strconv.Atoi(r.Find("span").Text())
		})

		entries = append(entries, entry)
	})

	return entries
}
