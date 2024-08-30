package adidas

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync/atomic"

	"vcrawler/internal/dto"
	"vcrawler/pkg/helpers"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

type scraper struct {
}

func (s *scraper) GetProductsURL(dumpLimit int) ([]string, error) {
	var (
		c                      *colly.Collector
		productURLs            []string
		currentPage, pageTotal int = 1, 0
	)

	c = helpers.GetCollector()

	c.OnRequest(func(r *colly.Request) {
		slog.Info(fmt.Sprintf("%d/%d: %s", currentPage, pageTotal, "visiting"), "url", r.URL.String())
	})

	// Handle the JSON response
	c.OnResponse(func(r *colly.Response) {
		var plr ProductListResponse
		// Unmarshal JSON into Go struct
		if err := json.Unmarshal(r.Body, &plr); err != nil {
			slog.Error("error at unmarshalling json", "error", err)
			return
		}

		productURLs = append(productURLs, plr.URLList()...)
		pageTotal = plr.SearchOptions.PageTotal
		currentPage = plr.CurrentPage()

		// Dump limit check
		if dumpLimit > 0 {
			dumpLimit -= len(productURLs)
			if dumpLimit <= 0 {
				return
			}
		}

		// Visit next page
		if plr.CurrentPage() < plr.SearchOptions.PageTotal {
			err := c.Visit(fmt.Sprintf(listingURLfmt, plr.CurrentPage()+1))
			if err != nil {
				slog.Error("error at visiting next page", "error", err)
			}
		}
	})

	// Handle request errors
	c.OnError(func(r *colly.Response, err error) {
		slog.Error("error at fetching:", "url", r.Request.URL.String(), "error", err)
	})

	// Start the request
	err := c.Visit(fmt.Sprintf(listingURLfmt, currentPage))
	if err != nil {
		return nil, err
	}

	// Wait until all asynchronous callbacks are complete
	c.Wait()

	return productURLs, nil
}

func (s *scraper) GetProductsDetail(productsURL []string) ([]dto.Product, error) {
	var (
		c         *colly.Collector
		products  []dto.Product
		totalURLs int64 = int64(len(productsURL)) // Total number of URLs
		completed int64                           // Counter for completed requests

	)

	c = helpers.GetCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", helpers.GetRandomUserAgent())
		slog.Info("visiting", "url", r.URL.String())
	})

	// Handle the JSON response
	c.OnResponse(func(r *colly.Response) {
		var pr ProductResponse
		// Unmarshal JSON into Go struct
		if err := json.Unmarshal(r.Body, &pr); err != nil {
			slog.Error("error at unmarshalling json", "error", err)
			return
		}

		product := pr.ToProduct(r.Request.URL.String())
		products = append(products, product)

		// Increment the counter and display progress
		completedCount := atomic.AddInt64(&completed, 1)
		percentage := float64(completedCount) / float64(totalURLs) * 100
		fmt.Printf("\rProgress: %.2f%% (%d/%d)\n", percentage, completedCount, totalURLs)

	})

	// Handle request errors
	c.OnError(func(r *colly.Response, err error) {
		slog.Error("error at fetching:", "url", r.Request.URL.String(), "error", err)

		// Still increment the counter for errors to avoid progress being stuck
		completedCount := atomic.AddInt64(&completed, 1)
		percentage := float64(completedCount) / float64(totalURLs) * 100
		fmt.Printf("\rProgress: %.2f%% (%d/%d)\n", percentage, completedCount, totalURLs)

	})

	// Set up a queue with only 1 consumer thread
	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})

	for _, url := range productsURL {
		productURL := url
		q.AddURL(productURL)
	}

	// Process the queue
	q.Run(c)

	// Wait until all asynchronous callbacks are complete
	c.Wait()

	return products, nil
}
