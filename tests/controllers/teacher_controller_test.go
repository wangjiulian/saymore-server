package controllers

import (
	"net/http"
	"testing"
	"time"
	
	"com.say.more.server/internal/app/controllers"
	"com.say.more.server/internal/app/models/form"
	"github.com/stretchr/testify/assert"
)

func TestGetTeacherList(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 注册路由 - 获取教师列表不需要认证
	r.GET("/teachers", controllers.GetTeacherList)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/teachers?page=1&size=10", nil)
	
	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTeacherDetail(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 注册路由 - 获取教师详情不需要认证
	r.GET("/teacher/:id", controllers.GetTeacherDetail)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/teacher/1", nil)
	
	// 检查响应 - 实际使用时需要先在数据库中准备测试数据
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTeacherAvailability(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 注册路由 - 获取教师可用时间不需要认证
	r.GET("/teacher/:id/availabilities", controllers.GetTeacherAvailabilities)
	
	// 当前时间的下一个周一
	now := time.Now()
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	nextMonday := now.AddDate(0, 0, daysUntilMonday)
	
	// 发送请求 - 获取下一周的可用时间
	startDate := uint64(nextMonday.Unix())
	endDate := uint64(nextMonday.AddDate(0, 0, 7).Unix())
	url := "/teacher/1/availabilities?start_date=" + time.Unix(int64(startDate), 0).Format("2006-01-02") + "&end_date=" + time.Unix(int64(endDate), 0).Format("2006-01-02")
	
	w := MakeRequest(t, r, http.MethodGet, url, nil)
	
	// 检查响应 - 实际使用时需要先在数据库中准备测试数据
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeacherLogin(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 注册路由
	r.POST("/teacher/login", controllers.TeacherLogin)
	
	// 准备请求体
	loginForm := form.LoginForm{
		Phone:    "13900139000",
		Password: "password123",
	}
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodPost, "/teacher/login", loginForm)
	
	// 检查响应 - 实际情况下应该使用模拟数据库
	// 此处只是结构性测试，预期会返回错误（因为没有真实数据）
	response := ParseResponse(t, w)
	assert.NotEqual(t, 0, response.Code) // 期望返回错误码
}

func TestGetTeacherProfile(t *testing.T) {
	// 设置路由
	r := SetupRouter()
	
	// 模拟教师认证
	MockTeacherJwtToken(r, 10)
	
	// 注册路由
	r.GET("/teacher", controllers.GetTeacher)
	
	// 发送请求
	w := MakeRequest(t, r, http.MethodGet, "/teacher", nil)
	
	// 检查响应 - 实际情况下应该使用模拟数据库
	assert.Equal(t, http.StatusOK, w.Code)
} 