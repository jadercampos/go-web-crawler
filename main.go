package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type movie struct {
	Title string
	Year  string
}
type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
}

func main() {
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()
	crawl(*month, *day)
}

func crawl(month int, day int) {
	c := colly.NewCollector(
		// using async makes you lose the sort order
		// colly.Async(true)
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	infoCollector := c.Clone()

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	infoCollector.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1.header > span.itemprop")
		tmpProfile.Photo = e.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = e.ChildText("#name-job-categories > a > span.itemprop")
		tmpProfile.BirthDate = e.ChildAttr("#name-born-info time", "datetime")

		tmpProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))
		/*profileDiv := e.DOM.Find("#name-bio-text > div.name-trivia-bio-text > div.inline")
		tmpProfile.Bio = profileDiv.Text()*/

		e.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = kf.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = kf.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})
		js, err := json.MarshalIndent(tmpProfile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Profile URL: ", r.URL.String())
	})

	// uncomment below line if you enable Async mode
	// c.Wait()
	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startUrl)
}
