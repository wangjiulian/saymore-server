package controllers

import (
	"com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/models/vo"
	"com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/types"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func StudentPackageList(c *gin.Context) {
	var (
		subjectId       int64
		subjectParentId int64
		teacherId       int64
	)

	subjectIdStr := strings.TrimSpace(c.Query("subject_id"))
	teacherIdStr := strings.TrimSpace(c.Query("teacher_id"))
	if subjectIdStr != "" {
		// Convert string to int
		subjectId, _ = strconv.ParseInt(subjectIdStr, 10, 64)
		if subjectId <= 0 {
			BadRequest(c, "Invalid subject_id")
			return
		}
	}
	if teacherIdStr != "" {
		// Convert string to int
		teacherId, _ = strconv.ParseInt(teacherIdStr, 10, 64)
		if subjectId <= 0 {
			BadRequest(c, "Invalid teacher_id")
			return
		}
	}

	ctx := types.From(c)
	if subjectId > 0 {
		// Query the parent subject of the current subject, if no parent exists, the current subject is the parent
		var subject db.Subject
		if err := repository.Repos.DB.Where(db.SubjectFields.ID, subjectId).First(&subject).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				NotFound(c, "Subject not found")
				return
			}
			InternalServerError(c, err.Error())
			return
		}
		if subject.ParentId > 0 {
			subjectParentId = int64(subject.ParentId)
		} else {
			subjectParentId = subjectId
		}
	}

	var packages []db.StudentPackage
	query := repository.Repos.DB.Where(db.StudentPackageFields.StudentId, ctx.StudentID()).Where(db.StudentPackageFields.LeftHours + " > 0")

	if subjectParentId > 0 {
		query = query.Where(db.StudentPackageFields.SubjectId, subjectParentId)
	}

	if teacherId > 0 {
		// When querying by teacher's subject, filter packages bound to the teacher
		query = query.Joins("left join "+db.StudentTeacherPackageTableName+" on "+db.StudentPackageFields.ID+"="+db.StudentTeacherPackageFields.StudentPackageId).Where(db.StudentTeacherPackageFields.TeacherId, teacherId)
	}

	if err := query.Find(&packages).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	var list []vo.StudentPackageList
	for _, p := range packages {
		list = append(list, vo.StudentPackageList{
			ID:        p.Id,
			StudentId: p.StudentId,
			Name:      p.Name,
			SubjectId: p.SubjectId,
			Hours:     p.Hours.String(),
			LeftHours: p.LeftHours.String(),
		})
	}

	SuccessWithData(c, list)
}

func StudentPackageDetail(c *gin.Context) {
	ctx := types.From(c)

	offset := (ctx.Page() - 1) * ctx.PageSize()

	studentPackageId := c.Param("student_package_id")
	var studentPackageDetails []db.StudentPackageDetail

	if err := repository.Repos.DB.Where(map[string]interface{}{
		db.StudentPackageDetailFields.StudentPackageId: studentPackageId,
		db.StudentPackageDetailFields.StudentId:        ctx.StudentID(),
	}).Offset(offset).Limit(ctx.PageSize()).Order(db.StudentPackageDetailFields.ID + " DESC").Find(&studentPackageDetails).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	var list []vo.StudentPackageDetail
	for _, p := range studentPackageDetails {
		list = append(list, vo.StudentPackageDetail{
			ID:         p.Id,
			StudentId:  p.StudentId,
			Name:       p.Title,
			Hours:      p.Hours.String(),
			LeftHours:  p.LeftHours.String(),
			Change:     p.Change,
			ChangeType: p.ChangeType,
			Time:       p.CreatedAt,
		})
	}

	result := vo.ListData{
		List: list,
		PageInfo: vo.PageInfo{
			Page:     ctx.Page(),
			PageSize: ctx.PageSize(),
		},
	}

	SuccessWithData(c, result)
}
