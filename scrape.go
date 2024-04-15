package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// struct Restuarant holds the data about a restaurant
type Restaurant struct {
	Name    string
	Cuisine string
	Rating  string // Michelin stars or Bib Gourmand
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// General selector for all the restuarants listed on the page
	restaurantContainerSelector := "section.section-main.search-results.search-listing-result div.restaurant__list-row.js-restaurant__list_items > div"

	// OnHTML callback
	c.OnHTML(restaurantContainerSelector, func(e *colly.HTMLElement) {
		// Extract restaurant data
		restaurant := Restaurant{
			Name:    e.ChildText("div > div.flex-fill > div > div:nth-child(1) > div:nth-child(2) > h3"),
			Cuisine: e.ChildText("div > div.flex-fill > div > div.row.flex-fill > div > div.align-items-end.js-match-height-bottom > div:nth-child(2)"),
			Rating:  extractRating(e.ChildAttrs("div > div.flex-fill > div > div:nth-child(1) > div:nth-child(1) > div > span > img", "src")),
		}

		if restaurant.Name != "" {
			// Print restuarant list with information
			fmt.Printf("Found restaurant: %s, Price range and Cuisine: %s, Rating: %s\n", restaurant.Name, restaurant.Cuisine, restaurant.Rating)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start scraping on the Micheline guide toronto page
	c.Visit("https://guide.michelin.com/ca/en/restaurants/1-star-michelin/2-stars-michelin/bib-gourmand?q=toronto%2C+Ontario%2C+Canada")
}

// extractRating finds the rating based on the img src
func extractRating(imageSrcs []string) string {
	starCount := 0
	for _, src := range imageSrcs {
		if strings.Contains(src, "1star") {
			starCount++
		} else if strings.Contains(src, "bib-gourmand") {
			return "Bib Gourmand "
		}
	}
	switch starCount {
	case 1:
		return "1 star"
	case 2:
		return "2 stars"
	default:
		return "No Rating"
	}
}
