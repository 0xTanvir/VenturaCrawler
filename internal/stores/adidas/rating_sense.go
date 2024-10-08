package adidas

import (
	"fmt"
	"log/slog"
	"strings"

	"vcrawler/internal/dto"
	"vcrawler/pkg/helpers"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	ratingSenseURL = "https://adidasjp.ugc.bazaarvoice.com/7896-ja_jp/%s/reviews.djs?format=embeddedhtml&productattribute_itemKcod=%s"
)

func GetRatingSense(articleCode, modelCode string) (string, string, []dto.RatingSense) {
	var (
		c               *colly.Collector
		ratingSenses    []dto.RatingSense
		rating          string
		recommendedRate string
	)

	c = helpers.GetCollector()
	// Handle the response
	c.OnResponse(func(r *colly.Response) {
		// Convert the response body to a string
		responseString := string(r.Body)

		cleanedString := cleanResponse(responseString)

		// Load the cleaned string into goquery for scraping
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(cleanedString))
		if err != nil {
			slog.Error("error at loading document", "error", err)
			return // Return to avoid further processing
		}

		// Get the rating value
		// example:  <span itemprop="ratingValue" class="BVRRNumber BVRRRatingNumber">4</span>
		rating = doc.Find("span[itemprop=ratingValue]").Text()

		// Get the recommanded rate
		// example: <span class="BVRRBuyAgainPercentage"> <span class="BVRRNumber">86%</span> </span>
		recommendedRate = doc.Find("span.BVRRBuyAgainPercentage span.BVRRNumber").Text()

		// Scrape the data
		doc.Find("div.BVRRRatingEntry").Each(func(i int, s *goquery.Selection) {
			ratingType := s.Find("div.BVRRRatingHeader").Text()
			ratingValue := s.Find("div.BVRRRatingRadioImage img").AttrOr("alt", "N/A")

			if ratingType != "" || ratingValue != "N/A" {
				ratingSense := dto.RatingSense{
					Type:  strings.TrimSpace(ratingType),
					Value: strings.TrimSpace(ratingValue),
				}

				ratingSenses = append(ratingSenses, ratingSense)
			}
		})
	})

	// Handle request errors
	c.OnError(func(r *colly.Response, err error) {
		slog.Error("error at fetching:", "url", r.Request.URL.String(), "error", err)
	})

	// Start the request
	err := c.Visit(fmt.Sprintf(ratingSenseURL, modelCode, articleCode))
	if err != nil {
		slog.Error("error at visiting rating sense", "error", err)
		return rating, recommendedRate, nil
	}

	// Wait until all asynchronous callbacks are complete
	c.Wait()

	return rating, recommendedRate, ratingSenses
}

// cleanResponse performs cleaning operations on the raw response string
func cleanResponse(response string) string {
	response = strings.Split(response, `"BVRRRatingSummarySourceID":"`)[1]
	response = strings.Split(response, `","BVRRSecondaryRatingSummarySourceID":"`)[0]
	response = strings.ReplaceAll(response, `\n`, "")
	response = strings.ReplaceAll(response, `\`, "")
	return response
}
