package vo

type (
	Dict struct {
		Subjects         []Subject `json:"subjects"`
		CourseCancelRule string    `json:"course_cancel_rule"`
	}

	Subject struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		ParentId   int    `json:"parent_id"`
		ChildCount int64  `json:"child_count"`
	}
)
