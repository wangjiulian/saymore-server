package form

type (
	BookRegularCourseForm struct {
		TeacherId              int64  `json:"teacher_id" form:"teacher_id"`
		TeacherAvailabilityIds string `json:"teacher_availability_ids" form:"teacher_availability_ids"`
		StudentPackageId       int64  `json:"student_package_id" form:"student_package_id"`
		SubjectId              int    `json:"subject_id" form:"subject_id"`
	}

	BookTrialCourseForm struct {
		TeacherId             int64 `json:"teacher_id" form:"teacher_id"`
		TeacherAvailabilityId int64 `json:"teacher_availability_id" form:"teacher_availability_id"`
		SubjectId             int   `json:"subject_id" form:"subject_id"`
	}

	CancelCourseForm struct {
		Reason string `json:"reason" form:"reason"`
	}

	CourseEvaluationForm struct {
		ContentQualityRating    int    `json:"content_quality_rating" form:"content_quality_rating"`
		InstructorClarityRating int    `json:"instructor_clarity_rating" form:"instructor_clarity_rating"`
		LearningGainRating      int    `json:"learning_gain_rating" form:"learning_gain_rating"`
		Content                 string `json:"content" form:"content"`
	}
)
