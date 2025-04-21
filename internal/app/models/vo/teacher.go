package vo

type Recommends struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Nickname          string `json:"nickname"`
	AvatarUrl         string `json:"avatar_url"`
	TeachingStartDate string `json:"teaching_start_date"`
	Gender            string `json:"gender"`
	EducationSchool   string `json:"education_school"`
	EducationLevel    string `json:"education_level"`
	CourseNum         string `json:"course_num"`
	SubjectIDs        []int  `json:"subject_ids"`
}

type TeacherDetail struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	Nickname             string   `json:"nickname"`
	AvatarUrl            string   `json:"avatar_url"`
	TeachingYears        string   `json:"teaching_years"`
	Gender               string   `json:"gender"`
	Background           string   `json:"background"`
	EducationSchool      string   `json:"education_school"`
	EducationLevel       string   `json:"education_level"`
	CourseNum            string   `json:"course_num"`
	TeachingExperience   string   `json:"teaching_experience"`
	Notes                string   `json:"notes"`
	TeachingAchievements string   `json:"teaching_achievements"`
	SuccessCases         string   `json:"success_cases"`
	SubjectIDs           []string `json:"subject_ids"`
	Evaluation           string   `json:"evaluation"`
}

type TeacherAvailability struct {
	ID        int64  `json:"id"`
	StartTime uint64 `json:"start_time"`
	EndTime   uint64 `json:"end_time"`
}
