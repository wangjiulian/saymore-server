package types

import (
	"com.say.more.server/internal/app/constant"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	defaultPageSize = 10
	defaultPage     = 1
)

// Context define some common params
type Context struct {
	context.Context
	page, pageSize int
	sortBy         string
	studentId      int64
	deviceType     string
}

// From create a new context
func From(ctx *gin.Context) *Context {
	page := ctx.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)

	pageSize := ctx.DefaultQuery("page_size", constant.PageSize)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	studentId := ctx.GetString("student_id")
	studentIdInt, _ := strconv.ParseInt(studentId, 10, 64)

	deviceType := ctx.GetString("device_type")

	return &Context{
		Context:    context.Background(),
		page:       pageInt,
		pageSize:   pageSizeInt,
		sortBy:     ctx.DefaultQuery("sort_by", ""),
		studentId:  studentIdInt,
		deviceType: deviceType,
	}
}

// DefaultContext create a default Context
func DefaultContext() *Context {
	return &Context{
		Context:  context.Background(),
		page:     defaultPage,
		pageSize: defaultPageSize,
	}
}

// WithPage set page info
func (ctx *Context) WithPage(page int) *Context {
	ctx.page = page
	return ctx
}

// Page get page number
func (ctx *Context) Page() int {
	return ctx.page
}

// WithPageSize set page info
func (ctx *Context) WithPageSize(pageSize int) *Context {
	ctx.pageSize = pageSize
	return ctx
}

// PageSize get page size
func (ctx *Context) PageSize() int {
	return ctx.pageSize
}

// WithSortBy set page info
func (ctx *Context) WithSortBy(sortBy string) *Context {
	ctx.sortBy = sortBy
	return ctx
}

// SortBy get sor by
func (ctx *Context) SortBy() string {
	return ctx.sortBy
}

// WithContext set a context.Context
func (ctx *Context) WithContext(c context.Context) *Context {
	ctx.Context = c
	return ctx
}

// WithUserID set a user_id
func (ctx *Context) WithStudentID(studentId int64) *Context {
	ctx.studentId = studentId
	return ctx
}

// UserID get user_id
func (ctx *Context) StudentID() int64 {
	return ctx.studentId
}

// WithDeviceType set a device type
func (ctx *Context) WithDeviceType(deviceType string) *Context {
	ctx.deviceType = deviceType
	return ctx
}

// DeviceType get device type
func (ctx *Context) DeviceType() string {
	return ctx.deviceType
}
