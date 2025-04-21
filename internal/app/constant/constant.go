package constant

const (
	AliOSSFileSizeLimit = 1024 * 1024 * 10 // 10MB

	CancelCourseDefaultInterval = 60                                                                              // Default time interval for no-fault course cancellation (minutes)
	CancelCourseDefaultRefund   = 0.5                                                                             // Default refund ratio for fault-based cancellations
	CancelCourseDefaultRule     = "Cancel more than 1 hour before class for a full refund; otherwise, 50% refund" // Default course cancellation rule
	CourseUnitDefault           = 0.5                                                                             // Default course unit ratio (50 minutes / 1 unit, 0.5 unit = 25 minutes)

	MaxFeedbackSize = 200 // Maximum length for feedback content

	PageSize = "10"
	// JwtTokenKey JWT token key
	JwtTokenKey = "say_more_secret"
	// redis key
	KeyToken          = "token:"
	DefaultDeviceType = "mini" // Default supported device type
)

// IgnoreAuthMethodMap Methods to ignore authentication
var IgnoreAuthMethodMap = map[string]string{
	"/api/v1/student/session-key":        "POST",
	"/api/v1/student/login":              "POST",
	"/api/v1/base/upload-file":           "POST",
	"/api/v1/base/dict":                  "GET",
	"/api/v1/teacher/search":             "GET",
	"/api/v1/teacher/:id/detail":         "GET",
	"/api/v1/teacher/:id/availabilities": "GET",
	"/api/v1/teacher/recommends":         "GET",
	"/api/v1/base/banners":               "GET",
	"/api/v1/course/evaluations":         "GET",
}
