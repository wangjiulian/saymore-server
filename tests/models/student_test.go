package models

import (
	"testing"
	"time"
	
	"com.say.more.server/internal/app/models/db"
	"github.com/stretchr/testify/assert"
)

func TestStudentStruct(t *testing.T) {
	// Create a new student instance
	student := db.Student{
		Base: db.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Phone:           "13800138000",
		AvatarUrl:       "https://example.com/avatar.jpg",
		Nickname:        "TestStudent",
		Gender:          1, // Male
		BirthDate:       time.Now().AddDate(-20, 0, 0), // 20 years ago
		StudentType:     2, // Elementary
		LearningPurpose: 3, // General English
		EnglishLevel:    4, // Fluent
		IsActive:        1, // Active
	}
	
	// Verify the values are set correctly
	assert.Equal(t, int64(1), student.ID)
	assert.Equal(t, "13800138000", student.Phone)
	assert.Equal(t, "TestStudent", student.Nickname)
	assert.Equal(t, 1, student.Gender)
	assert.Equal(t, 2, student.StudentType)
	assert.Equal(t, 3, student.LearningPurpose)
	assert.Equal(t, 4, student.EnglishLevel)
	assert.Equal(t, 1, student.IsActive)
	
	// Test TableName method
	assert.Equal(t, "students", student.TableName())
	
	// Test StudentFields constants
	assert.Equal(t, "students.id", db.StudentFields.ID)
	assert.Equal(t, "students.phone", db.StudentFields.Phone)
	assert.Equal(t, "students.nickname", db.StudentFields.Nickname)
} 