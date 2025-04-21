package db

const CourseEvaluationTableName = "course_evaluations"

type CourseEvaluation struct {
	Base
	TeacherId               int64  `gorm:"column:teacher_id;type:bigint(20);default:0;comment:老师ID;NOT NULL" json:"teacher_id"`
	StudentId               int64  `gorm:"column:student_id;type:bigint(20);default:0;comment:学生ID;NOT NULL" json:"student_id"`
	CourseId                int64  `gorm:"column:course_id;type:bigint(20);default:0;comment:课程ID;NOT NULL" json:"course_id"`
	ContentQualityRating    int    `gorm:"column:content_quality_rating;type:tinyint(4);default:3;comment:内容质量评分;NOT NULL" json:"content_quality_rating"`
	InstructorClarityRating int    `gorm:"column:instructor_clarity_rating;type:tinyint(4);default:3;comment:讲师清晰度评分;NOT NULL" json:"instructor_clarity_rating"`
	LearningGainRating      int    `gorm:"column:learning_gain_rating;type:tinyint(4);default:3;comment:学习收获评分;NOT NULL" json:"learning_gain_rating"`
	Content                 string `gorm:"column:content;type:varchar(255);comment:评价内容;NOT NULL" json:"content"`
}

func (CourseEvaluation) TableName() string {
	return CourseEvaluationTableName
}

var CourseEvaluationFields = struct {
	ID                      string
	TeacherId               string
	StudentId               string
	CourseId                string
	ContentQualityRating    string
	InstructorClarityRating string
	LearningGainRating      string
	Content                 string
	UpdatedAt               string
	CreatedAt               string
}{
	ID:                      CourseEvaluationTableName + ".id",
	TeacherId:               CourseEvaluationTableName + ".teacher_id",
	StudentId:               CourseEvaluationTableName + ".student_id",
	CourseId:                CourseEvaluationTableName + ".course_id",
	ContentQualityRating:    CourseEvaluationTableName + ".content_quality_rating",
	InstructorClarityRating: CourseEvaluationTableName + ".instructor_clarity_rating",
	LearningGainRating:      CourseEvaluationTableName + ".learning_gain_rating",
	Content:                 CourseEvaluationTableName + ".content",
	UpdatedAt:               CourseEvaluationTableName + ".updated_at",
	CreatedAt:               CourseEvaluationTableName + ".created_at",
}
