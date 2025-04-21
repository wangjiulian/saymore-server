package db

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	TeacherTableName = "teachers"
)

var (
	EducationLevelMap = map[int]string{
		1: "本科",
		2: "硕士",
		3: "博士",
	}

	GenderMap = map[int]string{
		1: "男",
		2: "女",
	}
)

// Teacher 老师表
type Teacher struct {
	Base
	Phone                string          `gorm:"column:phone;type:varchar(20);comment:手机号;NOT NULL" json:"phone"`
	Name                 string          `gorm:"column:name;type:varchar(50);comment:真实姓名;NOT NULL" json:"name"`
	Nickname             string          `gorm:"column:nickname;type:varchar(50);comment:昵称;NOT NULL" json:"nickname"`
	Gender               int             `gorm:"column:gender;type:tinyint(4);default:0;comment:性别：0-未知，1-男，2-女;NOT NULL" json:"gender"`
	CourseHours          decimal.Decimal `gorm:"column:course_hours;type:decimal(10,2);default:0;comment:授课量;NOT NULL" json:"course_num"`
	AvatarUrl            string          `gorm:"column:avatar_url;type:varchar(255);comment:头像URL;NOT NULL" json:"avatar_url"`
	Background           string          `gorm:"column:background;type:text;comment:个人背景" json:"background"`
	VideoUrl             string          `gorm:"column:video_url;type:varchar(255);comment:视频介绍URL;NOT NULL" json:"video_url"`
	EducationSchool      string          `gorm:"column:education_school;type:varchar(100);comment:毕业院校;NOT NULL" json:"education_school"`
	EducationLevel       int             `gorm:"column:education_level;type:tinyint(4);default:0;comment:学历水平：0-未知，1-本科，2-硕士，3-博士;NOT NULL" json:"education_level"`
	TeachingStartDate    time.Time       `gorm:"column:teaching_start_date;type:date;comment:授课开始时间;NOT NULL" json:"teaching_start_date"`
	Notes                string          `gorm:"column:notes;type:text;comment:备注信息" json:"notes"`
	TeachingExperience   string          `gorm:"column:teaching_experience;type:text;comment:教学经验案例" json:"teaching_experience"`
	SuccessCases         string          `gorm:"column:success_cases;type:text;comment:真实提分案例" json:"success_cases"`
	TeachingAchievements string          `gorm:"column:teaching_achievements;type:text;comment:教学成果" json:"teaching_achievements"`
	IsActive             int             `gorm:"column:is_active;type:tinyint(1);default:0;comment:是否启用;NOT NULL" json:"is_active"`
	IsRecommend          int             `gorm:"column:is_recommend;type:tinyint(1);default:0;comment:是否推荐;NOT NULL" json:"is_recommend"`
	Evaluation           decimal.Decimal `gorm:"column:evaluation;type:decimal(4,2);default:0.00;comment:评分;NOT NULL" json:"evaluation"`
}

func (Teacher) TableName() string {
	return TeacherTableName
}

var TeacherFields = struct {
	ID                   string
	CreatedAt            string
	UpdatedAt            string
	Phone                string
	Name                 string
	Nickname             string
	Gender               string
	CourseHours          string
	AvatarUrl            string
	Background           string
	VideoUrl             string
	EducationSchool      string
	EducationLevel       string
	TeachingStartDate    string
	Notes                string
	TeachingExperience   string
	SuccessCases         string
	TeachingAchievements string
	IsActive             string
	IsRecommend          string
	Evaluation           string
}{
	ID:                   TeacherTableName + ".id",
	CreatedAt:            TeacherTableName + ".created_at",
	UpdatedAt:            TeacherTableName + ".updated_at",
	Phone:                TeacherTableName + ".phone",
	Name:                 TeacherTableName + ".name",
	Nickname:             TeacherTableName + ".nickname",
	Gender:               TeacherTableName + ".gender",
	CourseHours:          TeacherTableName + ".course_hours",
	AvatarUrl:            TeacherTableName + ".avatar_url",
	Background:           TeacherTableName + ".background",
	VideoUrl:             TeacherTableName + ".video_url",
	EducationSchool:      TeacherTableName + ".education_school",
	EducationLevel:       TeacherTableName + ".education_level",
	TeachingStartDate:    TeacherTableName + ".teaching_start_date",
	Notes:                TeacherTableName + ".notes",
	TeachingExperience:   TeacherTableName + ".teaching_experience",
	SuccessCases:         TeacherTableName + ".success_cases",
	TeachingAchievements: TeacherTableName + ".teaching_achievements",
	IsActive:             TeacherTableName + ".is_active",
	IsRecommend:          TeacherTableName + ".is_recommend",
	Evaluation:           TeacherTableName + ".evaluation",
}
