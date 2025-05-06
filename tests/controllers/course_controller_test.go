package controllers

import (
	"net/http"
	"testing"
	"time"
	
	"com.say.more.server/internal/app/controllers"
	"com.say.more.server/internal/app/models/form"
	"github.com/stretchr/testify/assert"
)

func TestGetCourseList(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.GET("/courses", controllers.GetCourseList)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/courses?page=1&size=10", nil)
	
	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCourseDetail(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.GET("/course/:id", controllers.GetCourseDetail)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/course/1", nil)
	
	// 检查响应 - 实际使用时需要先在数据库中准备测试数据
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookCourse(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.POST("/course", controllers.BookCourse)
	
	// 准备请求体
	tomorrow := time.Now().AddDate(0, 0, 1)
	bookForm := form.CourseForm{
		TeacherId: 10,
		SubjectId: 2,
		Name:      "英语口语练习",
		StartTime: uint64(tomorrow.Unix()),
		EndTime:   uint64(tomorrow.Add(time.Hour).Unix()),
	}
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodPost, "/course", bookForm)
	
	// 这个测试不连接真实数据库，实际使用时需要模拟数据库
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCancelCourse(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.PUT("/course/:id/cancel", controllers.CancelCourse)
	
	// 准备请求体
	cancelForm := form.CourseCancel{
		Reason: "有事不能上课",
	}
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodPut, "/course/1/cancel", cancelForm)
	
	// 检查响应 - 实际使用时需要先在数据库中准备测试数据
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubmitCourseEvaluation(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.POST("/course/:id/evaluation", controllers.SubmitCourseEvaluation)
	
	// 准备请求体
	evalForm := form.CourseEvaluationForm{
		Score:    5,
		Content:  "老师讲课非常清晰，很满意这节课",
		IsPublic: 1,
	}
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodPost, "/course/1/evaluation", evalForm)
	
	// 检查响应 - 实际使用时需要先在数据库中准备测试数据
	assert.Equal(t, http.StatusOK, w.Code)
} 