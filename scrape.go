package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Restaurant struct {
	Name    string `json:"name"`
	Cuisine string `json:"cuisine"`
	Rating  string `json:"rating"`
}

type RestaurantList struct {
	Restaurants []Restaurant `json:"restaurants"`
}

func main() {
	c := colly.NewCollector()

	var restaurants RestaurantList

	restaurantContainerSelector := "section.section-main.search-results.search-listing-result div.restaurant__list-row.js-restaurant__list_items > div"

	c.OnHTML(restaurantContainerSelector, func(e *colly.HTMLElement) {
		restaurant := Restaurant{
			Name:    cleanText(e.ChildText("div > div.flex-fill > div > div:nth-child(1) > div:nth-child(2) > h3")),
			Cuisine: cleanText(e.ChildText("div > div.flex-fill > div > div.row.flex-fill > div > div.align-items-end.js-match-height-bottom > div:nth-child(2)")),
			Rating:  extractRating(e.ChildAttrs("div > div.flex-fill > div > div:nth-child(1) > div:nth-child(1) > div > span > img", "src")),
		}

		if restaurant.Name != "" {
			restaurants.Restaurants = append(restaurants.Restaurants, restaurant)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		jsonData, err := json.Marshal(restaurants)
		if err != nil {
			fmt.Println("Error serializing JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
		fmt.Println("Finished", r.Request.URL)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit("https://guide.michelin.com/ca/en/restaurants/1-star-michelin/2-stars-michelin/bib-gourmand?q=toronto%2C+Ontario%2C+Canada")
}

func extractRating(imageSrcs []string) string {
	starCount := 0
	for _, src := range imageSrcs {
		if strings.Contains(src, "1star") {
			starCount++
		} else if strings.Contains(src, "bib-gourmand") {
			return "Bib Gourmand"
		}
	}
	switch starCount {
	case 1:
		return "1 star"
	case 2:
		return "2 stars"
	case 3:
		return "3 stars"
	default:
		return "No Rating"
	}
}

func cleanText(input string) string {
	cleaned := strings.ReplaceAll(input, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\t", " ")
	// reg exp to replace multiple spaces with no space
	re := regexp.MustCompile(`\s+`)
	cleaned = re.ReplaceAllString(cleaned, "")
	return strings.TrimSpace(cleaned)
}
