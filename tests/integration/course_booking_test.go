package integration

import (
	"net/http"
	"testing"
	"time"
	
	"com.say.more.server/internal/app/controllers"
	"com.say.more.server/internal/app/models/form"
	"com.say.more.server/tests/controllers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCourseBookingFlow 测试预约课程完整流程
// 这个测试模拟了从预约课程到取消课程的完整流程
func TestCourseBookingFlow(t *testing.T) {
	// 跳过实际测试，因为我们没有设置测试数据库
	// 实际使用时，请移除此行并配置测试数据库
	t.Skip("Skipping integration test that requires test database")
	
	// 设置路由
	r := setupIntegrationRouter()
	
	// 步骤1: 预约课程
	courseId := bookNewCourse(t, r)
	
	// 步骤2: 查看课程详情
	verifyCourseDetails(t, r, courseId)
	
	// 步骤3: 取消课程
	cancelCourse(t, r, courseId)
	
	// 步骤4: 验证课程已取消
	verifyCourseIsCanceled(t, r, courseId)
}

// setupIntegrationRouter 设置集成测试路由
func setupIntegrationRouter() *gin.Engine {
	r := controllers.SetupRouter()
	
	// 模拟学生认证
	controllers.MockJwtToken(r, 123)
	
	// 注册路由
	r.POST("/course", controllers.BookCourse)
	r.GET("/course/:id", controllers.GetCourseDetail)
	r.PUT("/course/:id/cancel", controllers.CancelCourse)
	
	return r
}

// bookNewCourse 预约新课程
func bookNewCourse(t *testing.T, r *gin.Engine) int64 {
	// 准备请求体 - 预约明天的课程
	tomorrow := time.Now().AddDate(0, 0, 1)
	bookForm := form.CourseForm{
		TeacherId: 10,
		SubjectId: 2,
		Name:      "集成测试 - 英语口语课程",
		StartTime: uint64(tomorrow.Unix()),
		EndTime:   uint64(tomorrow.Add(time.Hour).Unix()),
	}
	
	// 发送请求
	w := controllers.MakeRequest(t, r, http.MethodPost, "/course", bookForm)
	
	// 检查响应
	response := controllers.AssertResponseSuccess(t, w)
	
	// 从响应中获取课程ID
	courseId, ok := response.Data.(map[string]interface{})["id"].(float64)
	assert.True(t, ok, "Response should contain course ID")
	
	return int64(courseId)
}

// verifyCourseDetails 验证课程详情
func verifyCourseDetails(t *testing.T, r *gin.Engine, courseId int64) {
	// 发送请求
	w := controllers.MakeRequest(t, r, http.MethodGet, "/course/"+string(courseId), nil)
	
	// 检查响应
	response := controllers.AssertResponseSuccess(t, w)
	
	// 验证课程信息
	course, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Response should contain course data")
	assert.Equal(t, float64(courseId), course["id"])
	assert.Equal(t, "集成测试 - 英语口语课程", course["name"])
	assert.Equal(t, float64(10), course["teacher_id"])
	assert.Equal(t, float64(2), course["subject_id"])
	assert.Equal(t, float64(1), course["status"]) // 待上课状态
}

// cancelCourse 取消课程
func cancelCourse(t *testing.T, r *gin.Engine, courseId int64) {
	// 准备请求体
	cancelForm := form.CourseCancel{
		Reason: "集成测试 - 取消课程",
	}
	
	// 发送请求
	w := controllers.MakeRequest(t, r, http.MethodPut, "/course/"+string(courseId)+"/cancel", cancelForm)
	
	// 检查响应
	controllers.AssertResponseSuccess(t, w)
}

// verifyCourseIsCanceled 验证课程已取消
func verifyCourseIsCanceled(t *testing.T, r *gin.Engine, courseId int64) {
	// 发送请求
	w := controllers.MakeRequest(t, r, http.MethodGet, "/course/"+string(courseId), nil)
	
	// 检查响应
	response := controllers.AssertResponseSuccess(t, w)
	
	// 验证课程已取消
	course, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Response should contain course data")
	assert.Equal(t, float64(4), course["status"]) // 已取消状态
	assert.Equal(t, "集成测试 - 取消课程", course["cancel_reason"])
} 