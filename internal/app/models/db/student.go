package db

import "time"

const StudentTableName = "students"

// Student table
type Student struct {
	Base
	Phone           string    `gorm:"column:phone;type:varchar(20);comment:Phone number;NOT NULL" json:"phone"`
	AvatarUrl       string    `gorm:"column:avatar_url;type:varchar(255);comment:Avatar URL;NOT NULL" json:"avatar_url"`
	Nickname        string    `gorm:"column:nickname;type:varchar(50);comment:Nickname;NOT NULL" json:"nickname"`
	Gender          int       `gorm:"column:gender;type:tinyint(1);default:0;comment:Gender: 1-Male, 2-Female, 0-Unspecified;NOT NULL" json:"gender"`
	BirthDate       time.Time `gorm:"column:birth_date;type:date;default:null;comment:Birth date" json:"birth_date"`
	StudentType     int       `gorm:"column:student_type;type:tinyint(1);default:0;comment:Student type: 1-Preschool, 2-Elementary, 3-Middle school, 4-High school, 5-College, 6-Professional;NOT NULL" json:"student_type"`
	LearningPurpose int       `gorm:"column:learning_purpose;type:tinyint(4);default:0;comment:Learning purpose: 1-IELTS, 2-TOEFL, 3-General English, 4-Business, 5-Social;NOT NULL" json:"learning_purpose"`
	EnglishLevel    int       `gorm:"column:english_level;type:tinyint(1);default:0;comment:English level: 1-Beginner, 2-Simple words, 3-Full sentences, 4-Fluent, 5-Excellent;NOT NULL" json:"english_level"`
	IsActive        int       `gorm:"column:is_active;type:tinyint(1);default:1;comment:Account status: 1-Active, 0-Disabled;NOT NULL" json:"is_active"`
}

func (Student) TableName() string {
	return StudentTableName
}

var StudentFields = struct {
	ID              string
	CreatedAt       string
	UpdatedAt       string
	Phone           string
	AvatarUrl       string
	Nickname        string
	Gender          string
	BirthDate       string
	StudentType     string
	LearningPurpose string
	EnglishLevel    string
	IsActive        string
}{
	ID:              StudentTableName + ".id",
	CreatedAt:       StudentTableName + ".created_at",
	UpdatedAt:       StudentTableName + ".updated_at",
	Phone:           StudentTableName + ".phone",
	AvatarUrl:       StudentTableName + ".avatar_url",
	Nickname:        StudentTableName + ".nickname",
	Gender:          StudentTableName + ".gender",
	BirthDate:       StudentTableName + ".birth_date",
	StudentType:     StudentTableName + ".student_type",
	LearningPurpose: StudentTableName + ".learning_purpose",
	EnglishLevel:    StudentTableName + ".english_level",
	IsActive:        StudentTableName + ".is_active",
}
