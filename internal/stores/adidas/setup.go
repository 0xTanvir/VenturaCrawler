package adidas

import "vcrawler/internal/definition"

const (
	baseURL       = "https://shop.adidas.jp"
	baseApiURLfmt = "https://shop.adidas.jp/f/v2/web/pub/products/article/%s/"
	listingURLfmt = "https://shop.adidas.jp/f/v1/pub/product/list?category=wear&gender=mens&limit=120&order=10&page=%d"
)

func GetAdidasStore() definition.Store {
	return &scraper{}
}
