package adidas

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"vcrawler/internal/dto"
)

// Price struct to hold price-related information
type Price struct {
	Badge struct {
		Color string `json:"color"`
		Text  string `json:"text"`
	} `json:"badge"`
	PriceDiscount       int    `json:"price_discount"`
	PriceDiscountStatus string `json:"price_discount_status"`
	PriceDiscountTax    string `json:"price_discount_tax"`
	PriceFixed          int    `json:"price_fixed"`
	PriceFixedTax       string `json:"price_fixed_tax"`
	PriceStatusMemo     string `json:"price_status_memo"`
}

// ImageNameData struct to hold image-related information
type ImageNameData struct {
	One   string `json:"1"`
	Two   string `json:"2,omitempty"`
	Hover string `json:"hover,omitempty"`
}

// Article struct to hold article-related information
type Article struct {
	Article             string `json:"article"`
	BrandName           string `json:"brand_name"`
	BrandSlug           string `json:"brand_slug"`
	BrandType           string `json:"brand_type"`
	CategoryNameEn      string `json:"category_name_en"`
	CategoryNameJa      string `json:"category_name_ja"`
	CategorySlug        string `json:"category_slug"`
	CategoryType        string `json:"category_type"`
	ColorName           string `json:"color_name"`
	ColorsCount         string `json:"colors_count"`
	DisplayBrand        string `json:"display_brand"`
	DisplayGender       string `json:"display_gender"`
	DisplayType         string `json:"display_type"`
	GenderName          string `json:"gender_name"`
	GenderSlug          string `json:"gender_slug"`
	GenderType          string `json:"gender_type"`
	GenreCode           string `json:"genre_code"`
	GenreName           string `json:"genre_name"`
	GroupName           string `json:"group_name"`
	GroupSlug           string `json:"group_slug"`
	ItemLimited         string `json:"item_limited"`
	ItemNameJa          string `json:"item_name_ja"`
	ItemStatus          string `json:"item_status"`
	LinkDetailPage      string `json:"link_detail_page"`
	MasterItemStatus    string `json:"master_item_status"`
	ModelCode           string `json:"model_code"`
	OptionalLabel       string `json:"optional_label"`
	OptionalLabelCode   string `json:"optional_label_code"`
	Price               Price  `json:"price"`
	PriceDiscount       int    `json:"price_discount"`
	PriceDiscountStatus string `json:"price_discount_status"`
	PriceDiscountTax    string `json:"price_discount_tax"`
	PriceFixed          int    `json:"price_fixed"`
	PriceFixedTax       string `json:"price_fixed_tax"`
	PriceStatusMemo     string `json:"price_status_memo"`
	ReleaseDate         string `json:"release_date"`
	ReviewCount         int    `json:"review_count"`
	ReviewRating        int    `json:"review_rating"`
	SportCode           string `json:"sport_code"`
	SportName           string `json:"sport_name"`
	SportSlug           string `json:"sport_slug"`
}

// Articles struct to hold multiple articles
type Articles map[string]Article

// SearchOptions struct to hold search options
type SearchOptions struct {
	Limit     string `json:"limit"`
	Offset    int    `json:"offset"`
	Page      string `json:"page"`
	PageTotal int    `json:"page_total"`
	Sort      string `json:"sort"`
}

// ProductListResponse struct to hold the full response
type ProductListResponse struct {
	Articles      Articles      `json:"articles"`
	SearchOptions SearchOptions `json:"search_options"`
}

func (plr ProductListResponse) URLList() []string {
	var urls []string
	for _, article := range plr.Articles {
		urls = append(urls, fmt.Sprintf(baseApiURLfmt, article.Article))
	}
	return urls
}

func (plr ProductListResponse) CurrentPage() int {
	pageNo, err := strconv.Atoi(plr.SearchOptions.Page)
	if err != nil {
		slog.Error("error at converting page number", "cause", err)
		return 0
	}

	return pageNo
}

func (pr ProductResponse) ToProduct() dto.Product {
	rating, recommendedRate, ratingSense := GetRatingSense(pr.Product.Article.ArticleCode, pr.Product.Model.ModelCode)

	return dto.Product{
		Name:        pr.Product.Article.Name,
		ModelCode:   pr.Product.Model.ModelCode,
		ArticleCode: pr.Product.Article.ArticleCode,
		Price: dto.Price{
			WithTax:      pr.Product.Article.Price.Current.WithTax,
			WithoutTax:   strconv.FormatFloat(pr.Product.Article.Price.Current.WithoutTax, 'f', 6, 64),
			DiscountType: pr.Product.Article.Price.DiscountType,
		},
		URL:             fmt.Sprintf("%s/products/%s", baseURL, pr.Product.Article.ArticleCode),
		Images:          pr.Images(),
		Breadcrumb:      pr.Breadcrumb(),
		Breadcrumbs:     pr.Breadcrumbs(),
		KWs:             pr.KWs(),
		Categories:      pr.Categories(),
		SizeChoice:      pr.SizeChoice(),
		Coordinates:     pr.Coordinates(),
		Description:     pr.Description(),
		Skus:            pr.Skus(),
		SizeCharts:      GetSizeCharts(pr.Product.Model.ModelCode),
		Technologies:    pr.Technologies(),
		ReviewCount:     fmt.Sprintf("%d", pr.Product.Model.Review.ReviewCount),
		Reviews:         pr.Reviews(),
		Rating:          rating,
		RecommendedRate: recommendedRate,
		RatingSenses:    ratingSense,
	}
}

func (pr ProductResponse) Reviews() []dto.Review {
	var result []dto.Review
	for _, review := range pr.Product.Model.Review.ReviewSeoLd {
		result = append(result, dto.Review{
			AuthorName:    review.Name,
			DatePublished: review.DatePublished,
			Body:          review.ReviewBody,
			BestRating:    fmt.Sprintf("%.1f", review.ReviewRating.BestRating),
			RatingValue:   review.ReviewRating.RatingValue,
		})
	}
	return result
}

func (pr ProductResponse) Technologies() []dto.Technology {
	var result []dto.Technology
	for _, tech := range pr.Product.Model.Description.Technology {
		result = append(result, dto.Technology{
			Name: tech.Name,
			Desc: tech.Text,
		})
	}
	return result
}

func (pr ProductResponse) Skus() []dto.Sku {
	var result []dto.Sku
	for _, sku := range pr.Product.Article.Skus {
		result = append(result, dto.Sku{
			SizeName: sku.SizeName,
			Status: dto.SkuStatus{
				IsStockEc:    sku.Status.InStockEc,
				IsStockStore: sku.Status.InStockStore,
				IsSoldOut:    sku.Status.IsSoldOut,
			},
		})
	}
	return result
}

func (pr ProductResponse) Description() dto.Description {
	return dto.Description{
		Title:   pr.Product.Article.Description.Messages.Title,
		General: pr.Product.Article.Description.Messages.MainText,
		Breads:  pr.Product.Article.Description.Messages.Breads,
	}
}

func (pr ProductResponse) Coordinates() []dto.Coordinate {
	var result []dto.Coordinate
	for _, article := range pr.Product.Article.Coordinates.Articles {
		result = append(result, dto.Coordinate{
			ProductName:  article.Name,
			ProductURL:   fmt.Sprintf("%s/products/%s", baseURL, article.ArticleCode),
			ProductImage: fmt.Sprintf("%s%s", baseURL, article.Image),
			ProductPrice: dto.Price{
				WithTax:      article.Price.Current.WithTax,
				WithoutTax:   strconv.FormatFloat(article.Price.Current.WithoutTax, 'f', 6, 64),
				DiscountType: article.Price.DiscountType,
			},
		})
	}
	return result
}

func (pr ProductResponse) SizeChoice() dto.SizeChoice {
	var availableSizes []string
	for _, sku := range pr.Product.Article.Skus {
		if sku.Status.InStockEc {
			availableSizes = append(availableSizes, sku.SizeName)
		}
	}

	return dto.SizeChoice{
		AvailableSize:  strings.Join(availableSizes, ", "),
		SenseOfTheSize: fmt.Sprintf("%.1f", pr.Product.Model.Review.FitbarScore),
	}
}

func (pr ProductResponse) KWs() string {
	var kws []string
	for _, category := range pr.Page.Categories {
		kws = append(kws, category.Label)
	}

	return strings.Join(kws, ", ")
}

func (pr ProductResponse) Categories() []dto.Category {
	var categories []dto.Category
	for _, category := range pr.Page.Categories {
		categories = append(categories, dto.Category{
			Label: category.Label,
			Link:  fmt.Sprintf("%s%s", baseURL, category.Link),
		})
	}
	return categories
}

func (pr ProductResponse) Breadcrumbs() []dto.Breadcrumb {
	var breadcrumbs []dto.Breadcrumb
	for _, bc := range pr.Page.Breadcrumbs {
		breadcrumbs = append(breadcrumbs, dto.Breadcrumb{
			Label:     bc.Label,
			SearchURL: fmt.Sprintf("%s%s", baseURL, bc.SearchUrl),
		})
	}
	return breadcrumbs
}

func (pr ProductResponse) Breadcrumb() string {
	var breadcrumbs []string
	for _, bc := range pr.Page.Breadcrumbs {
		breadcrumbs = append(breadcrumbs, bc.Label)
	}
	return strings.Join(breadcrumbs, " / ")
}

func (pr ProductResponse) Images() []string {
	var images []string
	for _, detail := range pr.Product.Article.Image.Details {
		images = append(images, fmt.Sprintf("%s%s", baseURL, detail.ImageUrl.Large))
	}
	return images
}

type ProductResponse struct {
	Page    Page    `json:"page"`
	Product Product `json:"product"`
}

type Page struct {
	Breadcrumbs []Breadcrumb `json:"breadcrumbs"`
	Categories  []Category   `json:"categories"`
}

type Breadcrumb struct {
	Label       string       `json:"label"`
	QueryParams []QueryParam `json:"queryParams"`
	SearchUrl   string       `json:"searchUrl"`
}

type QueryParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Category struct {
	GaEventLabel string `json:"gaEventLabel"`
	Label        string `json:"label"`
	Link         string `json:"link"`
}

type Product struct {
	Article ProductArticle `json:"article"`
	Model   Model          `json:"model"`
}

type ProductArticle struct {
	ArticleCode string       `json:"articleCode"`
	Coordinates Coordinates  `json:"coordinates"`
	Description Description  `json:"description"`
	Image       Image        `json:"image"`
	ModelCode   string       `json:"modelCode"`
	Name        string       `json:"name"`
	Price       ProductPrice `json:"price"`
	Skus        []Sku        `json:"skus"`
}

type Coordinates struct {
	Articles []ArticleItem `json:"articles"`
}

type ArticleItem struct {
	ArticleCode string       `json:"articleCode"`
	Image       string       `json:"image"`
	Name        string       `json:"name"`
	Price       ProductPrice `json:"price"`
}

type Description struct {
	Messages Messages `json:"messages"`
}

type Messages struct {
	Breads   []string `json:"breads"`
	MainText string   `json:"mainText"`
	Title    string   `json:"title"`
}

type Image struct {
	Details []ImageDetail `json:"details"`
}

type ImageDetail struct {
	Caption  string   `json:"caption"`
	ImageUrl ImageUrl `json:"imageUrl"`
}

type ImageUrl struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
}

type ProductPrice struct {
	Current      PriceDetail `json:"current"`
	DiscountType string      `json:"discountType"`
}

type PriceDetail struct {
	WithTax    string  `json:"withTax"`
	WithoutTax float64 `json:"withoutTax"`
}

type Sku struct {
	ArticleCode string `json:"articleCode"`
	SizeName    string `json:"sizeName"`
	Status      Status `json:"status"`
}

type AddToCartButton struct {
	Enabled bool   `json:"enabled"`
	Label   string `json:"label"`
}

type PurchaseInfo struct {
	Icon         string `json:"icon"`
	StockMessage string `json:"stockMessage"`
}

type Status struct {
	InStockEc    bool `json:"inStockEc"`
	InStockStore bool `json:"inStockStore"`
	IsSoldOut    bool `json:"isSoldOut"`
}

type Model struct {
	Description ModelDescription `json:"description"`
	ModelCode   string           `json:"modelCode"`
	Review      Review           `json:"review"`
}

type ModelDescription struct {
	Technology []Technology `json:"technology"`
}

type Technology struct {
	ImagePath string `json:"imagePath"`
	Name      string `json:"name"`
	Text      string `json:"text"`
}

type Review struct {
	FitbarScore float64       `json:"fitbarScore"`
	RatingAvg   float64       `json:"ratingAvg"`
	ReviewCount int           `json:"reviewCount"`
	ReviewSeoLd []ReviewSeoLd `json:"reviewSeoLd"`
}

type ReviewSeoLd struct {
	DatePublished string       `json:"datePublished"`
	Name          string       `json:"name"`
	ReviewBody    string       `json:"reviewBody"`
	ReviewRating  ReviewRating `json:"reviewRating"`
}

type ReviewRating struct {
	BestRating  float64 `json:"bestRating"`
	RatingValue string  `json:"ratingValue"`
}
