package letterboxdgo

import (
	"fmt"
	"log"
	"net/http"
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
	Title string
	TMDb  string
}

func GetFilm(slug string) *Film {

	f := new(Film)

	res, err := http.Get("https://letterboxd.com/film/" + slug)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Get title
	doc.Find("section.film-header-group span.name").Each(func(i int, s *goquery.Selection) {
		f.Title = s.Text()
	})

	// Get TMDb ID
	doc.Find("[data-track-action='TMDb']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		f.TMDb = strings.ReplaceAll(href, "https://www.themoviedb.org/movie/", "")
		f.TMDb = strings.ReplaceAll(f.TMDb, "/", "")
	})

	return f
}

func GetDiary(user string) []*DiaryEntry {
	// Get page count
	pageCount := 0

	res, err := http.Get(fmt.Sprintf("https://letterboxd.com/%s/films/diary/", user))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li.paginate-page").Each(func(i int, s *goquery.Selection) {
		pageCount = i
	})
	pageCount++

	var diary []*DiaryEntry
	for i := pageCount; i > 0; i-- {
		diary = append(diary, getDiaryPage(user, i)...)
	}

	// Sort entire diary by watch date
	sort.Slice(diary, func(i, j int) bool {
		return diary[i].Date.After(diary[j].Date)
	})

	return diary
}

func getDiaryPage(user string, page int) []*DiaryEntry {

	url := fmt.Sprintf("https://letterboxd.com/%s/films/diary/page/%v", user, page)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var entries []*DiaryEntry

	doc.Find("tr.diary-entry-row").Each(func(i int, s *goquery.Selection) {
		entry := new(DiaryEntry)

		// Get title and slug
		s.Find("td.td-film-details").Each(func(i int, r *goquery.Selection) {
			entry.Title = r.Find("a").Text()
			entry.Slug, _ = r.Find("div.poster").Attr("data-film-slug")
		})

		// Get watch date
		s.Find("td.diary-day a").Each(func(i int, r *goquery.Selection) {
			fullDatePath, _ := r.Attr("href")
			datePath := strings.ReplaceAll(fullDatePath, fmt.Sprintf("/%s/films/diary/for/", user), "")
			entry.Date, _ = time.Parse("2006/01/02/", datePath)
		})

		// Get rating
		s.Find("input.rateit-field").Each(func(i int, r *goquery.Selection) {
			sRating, _ := r.Attr("value")
			entry.Rating, _ = strconv.Atoi(sRating)
		})

		// Get Rewatch
		s.Find("td.td-rewatch").Each(func(i int, r *goquery.Selection) {
			classes, _ := r.Attr("class")
			if strings.Contains(classes, "icon-status-off") {
				entry.Rewatch = false
			} else {
				entry.Rewatch = true
			}
		})

		// Get Liked
		entry.Liked = false
		s.Find("td.td-like").Each(func(i int, r *goquery.Selection) {
			r.Find("span.hide-for-owner").Each(func(i int, l *goquery.Selection) {
				entry.Liked = true
			})
		})

		// Get release year
		s.Find("td.td-released").Each(func(i int, r *goquery.Selection) {
			entry.ReleaseYear, _ = strconv.Atoi(r.Find("span").Text())
		})

		entries = append(entries, entry)
	})

	return entries
}
