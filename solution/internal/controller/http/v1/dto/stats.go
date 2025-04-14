package dto

type Stats struct {
	ImpressionsCount int     `json:"impressions_count"`
	ClicksCount      int     `json:"clicks_count"`
	Conversion       float32 `json:"conversion"`
	SpentImpressions float32 `json:"spent_impressions"`
	SpentClicks      float32 `json:"spent_clicks"`
	SpentTotal       float32 `json:"spent_total"`
}

type DailyStats struct {
	Date             int     `json:"date"`
	ImpressionsCount int     `json:"impressions_count"`
	ClicksCount      int     `json:"clicks_count"`
	Conversion       float32 `json:"conversion"`
	SpentImpressions float32 `json:"spent_impressions"`
	SpentClicks      float32 `json:"spent_clicks"`
	SpentTotal       float32 `json:"spent_total"`
}
