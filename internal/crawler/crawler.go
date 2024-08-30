package crawler

import (
	"encoding/json"
	"log/slog"
	"os"

	"vcrawler/internal/definition"
)

type crawler struct {
}

func GetCrawler() definition.Crawler {
	return &crawler{}
}

func (c *crawler) Start(store definition.Store) error {
	dump := 200

	slog.Info("crawling products listing page")
	productsURL, err := store.GetProductsURL(dump)
	if err != nil {
		return err
	}

	productsURL = productsURL[:dump]
	products, err := store.GetProductsDetail(productsURL)
	if err != nil {
		return err
	}

	// Convert the products data to JSON format
	productsJSON, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	fileName := "products.json"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(productsJSON)
	if err != nil {
		return err
	}

	slog.Info("products data saved to", "file", fileName)
	return nil
}

func (c *crawler) Test(dump int, store definition.Store) error {
	slog.Info("crawling products listing page")
	productsURL, err := store.GetProductsURL(dump)
	if err != nil {
		return err
	}

	productsURL = productsURL[:dump]

	products, err := store.GetProductsDetail(productsURL)
	if err != nil {
		return err
	}

	for _, product := range products {
		s, err := json.MarshalIndent(product, "", "")
		if err != nil {
			slog.Error("error at marshalling product", "cause", err)
			continue
		}
		slog.Info("product detail", "product", string(s))
	}

	return nil
}
