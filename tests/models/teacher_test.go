package models

import (
	"testing"
	"time"
	
	"com.say.more.server/internal/app/models/db"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTeacherStruct(t *testing.T) {
	// Create decimal values for test
	courseHours := decimal.NewFromFloat(10.5)
	evaluation := decimal.NewFromFloat(4.5)
	
	// Create a teacher instance
	teacher := db.Teacher{
		Base: db.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Phone:                "13900139000",
		Name:                 "王老师",
		Nickname:             "Teacher Wang",
		Gender:               1, // 男
		CourseHours:          courseHours,
		AvatarUrl:            "https://example.com/teacher_avatar.jpg",
		Background:           "5年教学经验，专注英语口语教学",
		VideoUrl:             "https://example.com/intro_video.mp4",
		EducationSchool:      "北京大学",
		EducationLevel:       2, // 硕士
		TeachingStartDate:    time.Now().AddDate(-5, 0, 0), // 5 years ago
		Notes:                "擅长针对性教学",
		TeachingExperience:   "曾任职于多家知名英语培训机构",
		SuccessCases:         "学生英语水平提升显著，多人通过IELTS考试",
		TeachingAchievements: "获得年度优秀教师称号",
		IsActive:             1,
		IsRecommend:          1,
		Evaluation:           evaluation,
	}
	
	// Verify the values are set correctly
	assert.Equal(t, int64(1), teacher.ID)
	assert.Equal(t, "13900139000", teacher.Phone)
	assert.Equal(t, "王老师", teacher.Name)
	assert.Equal(t, "Teacher Wang", teacher.Nickname)
	assert.Equal(t, 1, teacher.Gender)
	assert.True(t, courseHours.Equal(teacher.CourseHours))
	assert.Equal(t, 2, teacher.EducationLevel)
	assert.Equal(t, 1, teacher.IsActive)
	assert.Equal(t, 1, teacher.IsRecommend)
	assert.True(t, evaluation.Equal(teacher.Evaluation))
	
	// Test TableName method
	assert.Equal(t, "teachers", teacher.TableName())
	
	// Test TeacherFields constants
	assert.Equal(t, "teachers.id", db.TeacherFields.ID)
	assert.Equal(t, "teachers.phone", db.TeacherFields.Phone)
	assert.Equal(t, "teachers.name", db.TeacherFields.Name)
	
	// Test constant maps
	assert.Equal(t, "本科", db.EducationLevelMap[1])
	assert.Equal(t, "硕士", db.EducationLevelMap[2])
	assert.Equal(t, "博士", db.EducationLevelMap[3])
	assert.Equal(t, "男", db.GenderMap[1])
	assert.Equal(t, "女", db.GenderMap[2])
} 