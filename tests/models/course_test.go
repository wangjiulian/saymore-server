package models

import (
	"testing"
	"time"
	
	"com.say.more.server/internal/app/models/db"
	"github.com/stretchr/testify/assert"
)

func TestCourseStruct(t *testing.T) {
	// Create a course instance
	course := db.Course{
		Base: db.Base{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CourseType:   db.CourseTypeRegular,
		Name:         "英语口语提高课程",
		SubjectId:    2,
		TeacherId:    10,
		StudentId:    20,
		Status:       db.CourseStatusWaiting,
		IsEvaluated:  db.CourseIsNotEvaluated,
		CancelReason: "",
		StartTime:    uint64(time.Now().Unix()),
		EndTime:      uint64(time.Now().Add(time.Hour).Unix()),
	}
	
	// Verify the values are set correctly
	assert.Equal(t, int64(1), course.ID)
	assert.Equal(t, db.CourseTypeRegular, course.CourseType)
	assert.Equal(t, "英语口语提高课程", course.Name)
	assert.Equal(t, 2, course.SubjectId)
	assert.Equal(t, int64(10), course.TeacherId)
	assert.Equal(t, int64(20), course.StudentId)
	assert.Equal(t, db.CourseStatusWaiting, course.Status)
	assert.Equal(t, db.CourseIsNotEvaluated, course.IsEvaluated)
	assert.Equal(t, "", course.CancelReason)
	
	// Test TableName method
	assert.Equal(t, "courses", course.TableName())
	
	// Test CourseFields constants
	assert.Equal(t, "courses.id", db.CourseFields.ID)
	assert.Equal(t, "courses.course_type", db.CourseFields.CourseType)
	assert.Equal(t, "courses.name", db.CourseFields.Name)
	assert.Equal(t, "courses.teacher_id", db.CourseFields.TeacherId)
	assert.Equal(t, "courses.student_id", db.CourseFields.StudentId)
	assert.Equal(t, "courses.status", db.CourseFields.Status)
	
	// Test course type constants
	assert.Equal(t, 1, db.CourseTypeRegular)
	assert.Equal(t, 2, db.CourseTypeTrial)
	
	// Test course status constants
	assert.Equal(t, 1, db.CourseStatusWaiting)
	assert.Equal(t, 2, db.CourseStatusOnGoing)
	assert.Equal(t, 3, db.CourseStatusFinished)
	assert.Equal(t, 4, db.CourseStatusCanceled)
	
	// Test course status map
	assert.Equal(t, "待上课", db.CourseStatusMap[db.CourseStatusWaiting])
	assert.Equal(t, "上课中", db.CourseStatusMap[db.CourseStatusOnGoing])
	assert.Equal(t, "已完成", db.CourseStatusMap[db.CourseStatusFinished])
	assert.Equal(t, "已取消", db.CourseStatusMap[db.CourseStatusCanceled])
	
	// Test course type map
	assert.Equal(t, "正式课", db.CourseTypeMap[db.CourseTypeRegular])
	assert.Equal(t, "试听课", db.CourseTypeMap[db.CourseTypeTrial])
} 