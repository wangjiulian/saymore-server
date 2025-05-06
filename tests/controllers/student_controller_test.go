package controllers

import (
	"net/http"
	"testing"
	
	"com.say.more.server/internal/app/controllers"
	"com.say.more.server/internal/app/models/form"
	"github.com/stretchr/testify/assert"
)

func TestGetStudentProfile(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.GET("/student", controllers.GetStudent)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/student", nil)
	
	// 如果数据库没有配置好，这里应该会返回错误
	// 这里只是测试路由和认证机制正常工作
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateStudentProfile(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.PUT("/student", controllers.UpdateStudent)
	
	// 准备请求体
	updateForm := form.StudentForm{
		Nickname:        "测试昵称",
		Gender:          1,
		StudentType:     2,
		LearningPurpose: 3,
		EnglishLevel:    4,
	}
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodPut, "/student", updateForm)
	
	// 这个测试不连接真实数据库，所以我们检查请求是否正常处理即可
	// 实际项目中，应该使用模拟数据库进行完整测试
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFeedback(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟学生认证
	MockJwtToken(r, 123)
	
	// 注册路由
	r.POST("/feedback", controllers.Feedback)
	
	// 准备请求体
	feedbackForm := form.FeedbackForm{
		Content: "这是一条测试反馈",
	}
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodPost, "/feedback", feedbackForm)
	
	// 同样，这个测试不连接真实数据库
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBanners(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 注册路由 - banner不需要认证
	r.GET("/banners", controllers.Banners)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/banners", nil)
	
	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)
} 