package vo

import "github.com/shopspring/decimal"

type CourseList struct {
	Id             int64  `json:"id"`
	TeacherId      int64  `json:"teacher_id"`
	SubjectId      int    `json:"subject_id"`
	Name           string `json:"name"`
	CourseTypeName string `json:"course_type_name"`
	TeacherName    string `json:"teacher_name"`
	Status         int    `json:"status"`
	IsEvaluated    int    `json:"is_evaluated"`
	StartTime      uint64 `json:"start_time"`
	EndTime        uint64 `json:"end_time"`
}

type CourseDetail struct {
	Id                     int64  `json:"id"`
	SubjectId              int    `json:"subject_id"`
	Name                   string `json:"name"`
	CourseTypeName         string `json:"course_type_name"`
	TeacherId              int64  `json:"teacher_id"`
	TeacherName            string `json:"teacher_name"`
	TeacherNickName        string `json:"teacher_nick_name"`
	TeacherEducationSchool string `json:"education_school"`
	TeacherEducationLevel  string `json:"education_level"`
	TeacherAvatar          string `json:"teacher_avatar"`
	Status                 int    `json:"status"`
	StatusText             string `json:"status_text"`
	StartTime              uint64 `json:"start_time"`
	EndTime                uint64 `json:"end_time"`
	IsEvaluated            int    `json:"is_evaluated"`
	CreatedTime            int64  `json:"created_time"`
}

type CourseEvaluationList struct {
	Id            int64           `json:"id"`
	StudentName   string          `json:"student_name"`
	StudentAvatar string          `json:"student_avatar"`
	AvgRating     decimal.Decimal `json:"avg_rating"`
	Content       string          `json:"content"`
	CreatedTime   int64           `json:"created_time"`
}
