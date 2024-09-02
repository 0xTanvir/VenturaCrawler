package dto

type ProductCsv struct {
	Name string `csv:"name" json:"name"`
	URL  string `csv:"url" json:"url"`
	// ModelCode                         string `csv:"model_code" json:"model_code"`
	// ArticleCode                       string `csv:"article_code" json:"article_code"`
	PriceWithTax                      string `csv:"price_with_tax" json:"price_with_tax"`
	PriceWithoutTax                   string `csv:"price_without_tax" json:"price_without_tax"`
	DiscountType                      string `csv:"discount_type" json:"discount_type"`
	Images                            string `csv:"images" json:"images"` // Concatenated string of image URLs
	Breadcrumb                        string `csv:"breadcrumb" json:"breadcrumb"`
	KWs                               string `csv:"kws" json:"kws"`
	SizeChoiceAvailable               string `csv:"available_size" json:"size_choice_available"`
	SizeChoiceSense                   string `csv:"sense_of_the_size" json:"size_choice_sense"`
	CoordinatesProductName            string `csv:"coordinates_product_name" json:"coordinates_product_name"`
	CoordinatesProductURL             string `csv:"coordinates_product_url" json:"coordinates_product_url"`
	CoordinatesProductImage           string `csv:"coordinates_product_image" json:"coordinates_product_image"`
	CoordinatesProductPriceWithTax    string `csv:"coordinates_product_price_with_tax" json:"coordinates_product_price_with_tax"`
	CoordinatesProductPriceWithoutTax string `csv:"coordinates_product_price_without_tax" json:"coordinates_product_price_without_tax"`
	CoordinatesProductDiscountType    string `csv:"coordinates_product_discount_type" json:"coordinates_product_discount_type"`
	DescriptionTitle                  string `csv:"description_title" json:"description_title"`
	DescriptionGeneral                string `csv:"general_description" json:"description_general"`
	DescriptionBreads                 string `csv:"general_itemization_description" json:"description_breads"` // Concatenated string of breads
	SizeCharts                        string `csv:"size_charts" json:"size_charts"`                            // Concatenated string of sizes and measurements
	Technologies                      string `csv:"special_function" json:"technologies"`                      // Concatenated string of technologies
	ReviewCount                       string `csv:"review_count" json:"review_count"`
	Reviews                           string `csv:"reviews" json:"reviews"` // Concatenated string of reviews
	Rating                            string `csv:"rating" json:"rating"`
	RecommendedRate                   string `csv:"recommended_rate" json:"recommended_rate"`
	RatingSenses                      string `csv:"rating_senses" json:"rating_senses"` // Concatenated string of rating senses
}
