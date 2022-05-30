package main

import (
	"flag"

	"github.com/jadercampos/go-web-crawler/controllers"
)

func main() {
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()
	controllers.CrawlMovie(*month, *day)
}
