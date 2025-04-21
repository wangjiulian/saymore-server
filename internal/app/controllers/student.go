package controllers

import (
	"com.say.more.server/internal/app/constant"
	"com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/models/form"
	"com.say.more.server/internal/app/models/vo"
	"com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/types"
	"com.say.more.server/internal/pkg/token"
	"com.say.more.server/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

func StudentSessionKey(c *gin.Context) {
	params := &form.CodeForm{}
	err := c.ShouldBind(params)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	params.Code = strings.TrimSpace(params.Code)
	if params.Code == "" {
		BadRequest(c, "code cannot be empty")
		return
	}

	session, err := repository.Repos.WeChat.GetSessionKey(params.Code)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	SuccessWithData(c, session.SessionKey)
}

func StudentLogin(c *gin.Context) {
	params := &form.LoginForm{}
	err := c.ShouldBind(params)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	params.EncryptedData = strings.TrimSpace(params.EncryptedData)
	params.Iv = strings.TrimSpace(params.Iv)
	params.SessionKey = strings.TrimSpace(params.SessionKey)
	if params.EncryptedData == "" || params.Iv == "" || params.SessionKey == "" {
		BadRequest(c, "invalid params")
		return
	}

	phone, err := repository.Repos.WeChat.DecryptWeChatPhoneData(params.EncryptedData, params.Iv, params.SessionKey)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	var student db.Student
	if err := repository.Repos.DB.Where(db.StudentFields.Phone, phone).First(&student).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			InternalServerError(c, err.Error())
			return
		}
	}

	if student.Id <= 0 {
		// Register
		student.Phone = phone
		// Take the last four digits
		student.Nickname = phone[len(phone)-4:] + " user"
		if err := repository.Repos.DB.Create(&student).Error; err != nil {
			InternalServerError(c, err.Error())
			return
		}
	}

	// Generate token
	token, err := token.RefreshToken(student.Id, student.Phone, constant.DefaultDeviceType)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	SuccessWithData(c, vo.LoginVo{
		StudentId: student.Id,
		Token:     token,
	})
}

func StudentDetail(c *gin.Context) {
	ctx := types.From(c)
	var student db.Student
	if err := repository.Repos.DB.Where(db.StudentFields.ID, ctx.StudentID()).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "student not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	var birthDate = student.BirthDate.Format("2006-01-02")
	if birthDate == "0001-01-01" {
		birthDate = ""
	}

	studentInfo := vo.StudentInfoVo{
		ID:              student.Id,
		Phone:           utils.MaskPhone(student.Phone),
		AvatarUrl:       student.AvatarUrl,
		Nickname:        student.Nickname,
		Gender:          student.Gender,
		BirthDate:       birthDate,
		StudentType:     student.StudentType,
		LearningPurpose: student.LearningPurpose,
		EnglishLevel:    student.EnglishLevel,
	}

	SuccessWithData(c, studentInfo)
}

func StudentEdit(c *gin.Context) {
	ctx := types.From(c)
	params := &form.StudentEditForm{}
	err := c.ShouldBind(params)
	if err != nil {
		BadRequest(c, "invalid params")
		return
	}

	// TODO: Validate data
	birthDate, err := time.Parse("2006-01-02", params.BirthDate)
	if err != nil {
		BadRequest(c, "invalid birth_date")
		return
	}
	upStudent := repository.Repos.DB.Model(&db.Student{}).Where(db.StudentFields.ID, ctx.StudentID()).Updates(db.Student{
		AvatarUrl:       params.Avatar,
		Nickname:        params.Nickname,
		Gender:          params.Gender,
		BirthDate:       birthDate,
		StudentType:     params.StudentType,
		LearningPurpose: params.LearningPurpose,
		EnglishLevel:    params.EnglishLevel,
	})
	if upStudent.Error != nil {
		InternalServerError(c, upStudent.Error.Error())
		return
	}

	Success(c)
}

func StudentTrial(c *gin.Context) {
	ctx := types.From(c)
	var trialCount int64
	if err := repository.Repos.DB.Table(db.StudentTrialQuota{}.TableName()).Where(map[string]interface{}{db.StudentTrialQuotaFields.StudentId: ctx.StudentID(), db.StudentTrialQuotaFields.CourseId: 0}).
		Count(&trialCount).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithData(c, trialCount)
}
