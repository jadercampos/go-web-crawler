package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/jadercampos/go-web-crawler/models"
)

func CrawlGames() {
	c := colly.NewCollector(colly.AllowedDomains("imdb.com", "www.imdb.com"))

	infoCollector := c.Clone()
	tmpProfile := models.Game{}
	c.OnHTML(".lister-item", func(e *colly.HTMLElement) {
		tmpProfile.Img = e.ChildAttr("img.loadlate", "src")
		tmpProfile.Description = e.ChildText("div.lister-item-content > h3.lister-item-header > a")
		profileUrl := e.ChildAttr("div.lister-item-content > h3.lister-item-header > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	infoCollector.OnHTML("main", func(e *colly.HTMLElement) {
		tmpProfile.Title = e.ChildText(".sc-b73cd867-0")
		e.ForEach("a.sc-16ede01-3", func(_ int, kf *colly.HTMLElement) {
			tmpProfile.Categories = append(tmpProfile.Categories, kf.ChildText("ul > li"))
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

	c.Visit("https://www.imdb.com/search/title/?sort=user_rating,desc&title_type=game")
}
