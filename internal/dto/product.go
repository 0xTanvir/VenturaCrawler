package dto

import "strings"

type SizeChoice struct {
	AvailableSize  string `csv:"available_size" json:"available_size"`
	SenseOfTheSize string `csv:"sense_of_the_size" json:"sense_of_the_size"`
}

type Price struct {
	WithTax      string `csv:"with_tax" json:"with_tax"`
	WithoutTax   string `csv:"without_tax" json:"without_tax"`
	DiscountType string `csv:"discount_type" json:"discount_type"`
}

type Category struct {
	Label string `csv:"label" json:"label"`
	Link  string `csv:"link" json:"link"`
}

type Breadcrumb struct {
	Label     string `csv:"label" json:"label"`
	SearchURL string `csv:"search_url" json:"search_url"`
}

type Coordinate struct {
	ProductName  string `csv:"product_name" json:"product_name"`
	ProductURL   string `csv:"product_url" json:"product_url"`
	ProductImage string `csv:"product_image" json:"product_image"`
	ProductPrice Price  `csv:"product_price" json:"product_price"`
}

type Description struct {
	Title   string   `csv:"title" json:"title"`
	General string   `csv:"general" json:"general"`
	Breads  []string `json:"breads"`
}

type SkuStatus struct {
	IsStockEc    bool `csv:"is_stock" json:"is_stock"`
	IsStockStore bool `csv:"is_stock_store" json:"is_stock_store"`
	IsSoldOut    bool `csv:"is_sold_out" json:"is_sold_out"`
}

type Sku struct {
	SizeName string    `csv:"size_name" json:"size_name"`
	Code     string    `csv:"code" json:"code"`
	Status   SkuStatus `csv:"status" json:"status"`
}

type Measurement struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SizeChart struct {
	Size         string        `csv:"size" json:"size"`
	Measurements []Measurement `json:"measurements"`
}

type Technology struct {
	Name string `csv:"name" json:"name"`
	Desc string `csv:"desc" json:"desc"`
}

type Review struct {
	AuthorName    string `csv:"author_name" json:"author_name"`
	DatePublished string `csv:"date_published" json:"date_published"`
	Body          string `csv:"body" json:"body"`
	BestRating    string `csv:"best_rating" json:"best_rating"`
	RatingValue   string `csv:"rating_value" json:"rating_value"`
}

type RatingSense struct {
	Type  string `csv:"type" json:"type"`
	Value string `csv:"value" json:"value"`
}

type Product struct {
	Name            string        `csv:"name" json:"name"`
	ModelCode       string        `csv:"model_code" json:"model_code"`
	ArticleCode     string        `csv:"article_code" json:"article_code"`
	Price           Price         `csv:"price" json:"price"`
	URL             string        `csv:"url" json:"url"`
	Images          []string      `csv:"images" json:"images"`
	Breadcrumb      string        `json:"breadcrumb"`
	Breadcrumbs     []Breadcrumb  `json:"breadcrumbs"`
	KWs             string        `csv:"kws" json:"kws"`
	Categories      []Category    `json:"categories"`
	SizeChoice      SizeChoice    `json:"size_choice"`
	Coordinates     []Coordinate  `json:"coordinates"`
	Description     Description   `json:"description"`
	Skus            []Sku         `json:"skus"`
	SizeCharts      []SizeChart   `json:"size_charts"`
	Technologies    []Technology  `json:"technologies,omitempty"`
	ReviewCount     string        `json:"review_count"`
	Reviews         []Review      `json:"reviews"`
	Rating          string        `json:"rating"`
	RecommendedRate string        `json:"recommended_rate"`
	RatingSenses    []RatingSense `json:"rating_senses"`
}

func (p Product) ToCsv() ProductCsv {
	// Concatenate image URLs
	images := strings.Join(p.Images, "; ")

	// Concatenate description breads
	descriptionBreads := strings.Join(p.Description.Breads, "; ") // General Itemization Description

	// Concatenate sizes and measurements
	var sizeChartDetails []string
	for _, sizeChart := range p.SizeCharts {
		var measurements []string
		for _, measurement := range sizeChart.Measurements {
			measurements = append(measurements, measurement.Type+": "+measurement.Value)
		}
		sizeChartDetails = append(sizeChartDetails, sizeChart.Size+" ("+strings.Join(measurements, ", ")+")")
	}
	sizeCharts := strings.Join(sizeChartDetails, "; ")

	// Concatenate technologies
	var technologies []string
	for _, tech := range p.Technologies {
		technologies = append(technologies, tech.Name+": "+tech.Desc)
	}
	technologiesStr := strings.Join(technologies, "; ")

	// Concatenate reviews
	var reviews []string
	for _, review := range p.Reviews {
		reviews = append(reviews, review.AuthorName+" ("+review.DatePublished+"): "+review.Body+" ["+review.RatingValue+"/"+review.BestRating+"]")
	}
	reviewsStr := strings.Join(reviews, "; ")

	// Concatenate rating senses
	var ratingSenses []string
	for _, ratingSense := range p.RatingSenses {
		ratingSenses = append(ratingSenses, ratingSense.Type+": "+ratingSense.Value)
	}
	ratingSensesStr := strings.Join(ratingSenses, "; ")

	// Convert to ProductCsv
	pCsv := ProductCsv{
		Name:                p.Name,
		PriceWithTax:        p.Price.WithTax,
		PriceWithoutTax:     p.Price.WithoutTax,
		DiscountType:        p.Price.DiscountType,
		URL:                 p.URL,
		Images:              images,
		Breadcrumb:          p.Breadcrumb,
		KWs:                 p.KWs,
		SizeChoiceAvailable: p.SizeChoice.AvailableSize,
		SizeChoiceSense:     p.SizeChoice.SenseOfTheSize,
		DescriptionTitle:    p.Description.Title,
		DescriptionGeneral:  p.Description.General,
		DescriptionBreads:   descriptionBreads,
		SizeCharts:          sizeCharts,
		Technologies:        technologiesStr,
		ReviewCount:         p.ReviewCount,
		Reviews:             reviewsStr,
		Rating:              p.Rating,
		RecommendedRate:     p.RecommendedRate,
		RatingSenses:        ratingSensesStr,
	}

	if len(p.Coordinates) > 0 {
		pCsv.CoordinatesProductName = p.Coordinates[0].ProductName
		pCsv.CoordinatesProductURL = p.Coordinates[0].ProductURL
		pCsv.CoordinatesProductImage = p.Coordinates[0].ProductImage
		pCsv.CoordinatesProductPriceWithTax = p.Coordinates[0].ProductPrice.WithTax
		pCsv.CoordinatesProductPriceWithoutTax = p.Coordinates[0].ProductPrice.WithoutTax
		pCsv.CoordinatesProductDiscountType = p.Coordinates[0].ProductPrice.DiscountType
	}

	return pCsv
}
