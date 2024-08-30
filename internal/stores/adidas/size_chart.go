package adidas

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"vcrawler/internal/dto"
	"vcrawler/pkg/helpers"

	"github.com/gocolly/colly/v2"
)

const (
	sizeChartURL = "https://shop.adidas.jp/f/v1/pub/size_chart/%s"
)

func GetSizeCharts(modelCode string) []dto.SizeChart {
	var (
		c          *colly.Collector
		sizeCharts []dto.SizeChart
	)

	c = helpers.GetCollector()

	// Handle the JSON response
	c.OnResponse(func(r *colly.Response) {
		var scr SizeChartResponse
		// Unmarshal JSON into Go struct
		if err := json.Unmarshal(r.Body, &scr); err != nil {
			slog.Error("error at unmarshalling json", "error", err)
			return
		}

		sizeCharts = scr.getSizeCharts()
	})

	// Handle request errors
	c.OnError(func(r *colly.Response, err error) {
		slog.Error("error at fetching:", "url", r.Request.URL.String(), "error", err)
	})

	// Start the request
	err := c.Visit(fmt.Sprintf(sizeChartURL, modelCode))
	if err != nil {
		slog.Error("error at visiting size chart", "error", err)
		return nil
	}

	// Wait until all asynchronous callbacks are complete
	c.Wait()

	return sizeCharts
}

func (scr SizeChartResponse) getSizeCharts() []dto.SizeChart {
	var sizeCharts []dto.SizeChart

	for _, sizeData := range scr.SizeChart {
		for _, body := range sizeData.Body {
			var measurements []dto.Measurement

			for bi, v := range body {
				m := dto.Measurement{
					Type:  scr.SizeChart["0"].Header["0"][bi].Value,
					Value: v.Value,
				}

				if m.Type != "" {
					measurements = append(measurements, m)
				}
			}

			sizeChart := dto.SizeChart{
				Size:         body["0"].Value,
				Measurements: measurements,
			}
			sizeCharts = append(sizeCharts, sizeChart)
		}
	}

	return sizeCharts
}

type SizeImage struct {
	Path string `json:"path"`
}

type SizeValue struct {
	Value string `json:"value"`
}

type SizeHeader struct {
	Value string `json:"value"`
}

type SizeBody struct {
	Value string `json:"value"`
}

type Body struct {
	XS   SizeBody `json:"0"`
	S    SizeBody `json:"1"`
	M    SizeBody `json:"2"`
	L    SizeBody `json:"3"`
	XL   SizeBody `json:"4"`
	XXL  SizeBody `json:"5"`
	XXXL SizeBody `json:"6"`
}

type SizeChartData struct {
	Body          map[string]map[string]SizeValue  `json:"body"`
	HasActualSize bool                             `json:"has_actual_size"`
	Header        map[string]map[string]SizeHeader `json:"header"`
}

type SizeChartResponse struct {
	IsExactFlag int                      `json:"is_exact_flag"`
	SizeChart   map[string]SizeChartData `json:"size_chart"`
}
