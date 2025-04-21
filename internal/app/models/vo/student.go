package vo

type LoginVo struct {
	StudentId int64  `json:"student_id"`
	Token     string `json:"token"`
}

type StudentInfoVo struct {
	ID              int64  `json:"id"`
	Phone           string `json:"phone"`
	AvatarUrl       string `json:"avatar"`
	Nickname        string `json:"nickname"`
	Gender          int    `json:"gender"`
	BirthDate       string `json:"birth_date"`
	StudentType     int    `json:"student_type"`
	LearningPurpose int    `json:"learning_purpose"`
	EnglishLevel    int    `json:"english_level"`
}
