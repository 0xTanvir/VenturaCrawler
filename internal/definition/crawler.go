package definition

import (
	"vcrawler/internal/dto"
)

type Downloader interface {
	// DownloadImage(url string) (imageData []byte, err error)
}

type Store interface {
	// GetProductsURL returns a list of product URLs from the listing page
	GetProductsURL(dumpLimit int) ([]string, error)
	// GetProductDetail returns the product details from the product page
	GetProductsDetail(productsURL []string) ([]dto.Product, error)

	// Downloader implements the downloader for the store
	Downloader
}

type Crawler interface {
	Start(store Store) error
	Test(dumpLimit int, store Store) error
}
