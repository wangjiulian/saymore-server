package form

type (
	FeedbackForm struct {
		Content string `json:"content" form:"content"`
	}
)
