package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type SearchTeacher struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Nickname          string          `json:"nickname"`
	AvatarUrl         string          `json:"avatar_url"`
	TeachingStartDate time.Time       `json:"teaching_start_date"`
	Gender            int             `json:"gender"`
	EducationSchool   string          `json:"education_school"`
	EducationLevel    int             `json:"education_level"`
	CourseHours       decimal.Decimal `json:"course_hours"`
}
