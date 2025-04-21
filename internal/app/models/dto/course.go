package dto

import "github.com/shopspring/decimal"

type CourseAvgRatings struct {
	AvgContentQuality decimal.Decimal
	AvgClarity        decimal.Decimal
	AvgLearningGain   decimal.Decimal
}
