package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"com.say.more.server/internal/app/models/vo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupRouter 设置测试路由
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// MakeRequest 发送测试请求并返回响应
func MakeRequest(t *testing.T, r *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		assert.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	req, err := http.NewRequest(method, url, reqBody)
	assert.NoError(t, err)
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ParseResponse 解析响应体
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder) *vo.Response {
	var response vo.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return &response
}

// AssertResponseSuccess 断言响应成功
func AssertResponseSuccess(t *testing.T, w *httptest.ResponseRecorder) *vo.Response {
	assert.Equal(t, http.StatusOK, w.Code)
	response := ParseResponse(t, w)
	assert.Equal(t, 0, response.Code)
	return response
}

// AssertResponseError 断言响应错误
func AssertResponseError(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) *vo.Response {
	response := ParseResponse(t, w)
	assert.Equal(t, expectedCode, response.Code)
	return response
}

// MockJwtToken 模拟JWT令牌
func MockJwtToken(r *gin.Engine, studentId int64) {
	r.Use(func(c *gin.Context) {
		c.Set("student_id", studentId)
		c.Next()
	})
}

// MockTeacherJwtToken 模拟教师JWT令牌
func MockTeacherJwtToken(r *gin.Engine, teacherId int64) {
	r.Use(func(c *gin.Context) {
		c.Set("teacher_id", teacherId)
		c.Next()
	})
} 