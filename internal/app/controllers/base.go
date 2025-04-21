package controllers

import (
	"com.say.more.server/internal/app/constant"
	"com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/models/form"
	"com.say.more.server/internal/app/models/vo"
	"com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
	"time"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		BadRequest(c, "invalid file")
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	defer src.Close()
	fileExt := filepath.Ext(fileHeader.Filename)
	if fileHeader.Size > constant.AliOSSFileSizeLimit {
		BadRequest(c, "file too large")
		return
	}

	aliOss := repository.Repos.AliOss
	key := fmt.Sprintf("uploads/%s", uuid.Must(uuid.NewRandom()).String()+fileExt)
	if err := aliOss.PutObject(key, src, nil); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithData(c, aliOss.Get(key))
}

func Dict(c *gin.Context) {

	subjectIdName, err := getDictSubject()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	dict := vo.Dict{
		Subjects:         subjectIdName,
		CourseCancelRule: repository.Repos.Config.Course.CancelRule,
	}

	SuccessWithData(c, dict)
}

func Feedback(c *gin.Context) {
	currentTime := time.Now().Unix()
	ctx := types.From(c)
	var params form.FeedbackForm
	if err := c.ShouldBind(&params); err != nil {
		BadRequest(c, "Invalid params")
		return
	}
	params.Content = strings.TrimSpace(params.Content)
	if params.Content == "" {
		BadRequest(c, "content cannot be empty")
		return
	}
	if len(params.Content) > constant.MaxFeedbackSize {
		BadRequest(c, "content too long")
		return
	}

	if err := repository.Repos.DB.Create(&db.Feedbacks{
		StudentId: ctx.StudentID(),
		Content:   params.Content,
		Base: db.Base{
			UpdatedAt: currentTime,
			CreatedAt: currentTime,
		},
	}).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c)
}

func Banners(c *gin.Context) {

	var banners []db.Banner
	if err := repository.Repos.DB.Find(&banners).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessWithData(c, banners)
}

func getDictSubject() ([]vo.Subject, error) {
	tx := repository.Repos.DB

	// Query the subject table
	var subjects []db.Subject
	if err := tx.Order(db.SubjectFields.ID + " ASC, " + db.SubjectFields.SortOrder + " ASC").Find(&subjects).Error; err != nil {
		return nil, err
	}
	subjectIdNameMap := make(map[int64]vo.Subject)
	for _, subject := range subjects {
		subjectIdNameMap[subject.Id] = vo.Subject{
			Name: subject.Name,
		}
		if subject.ParentId > 0 {
			parent := subjectIdNameMap[int64(subject.ParentId)]
			parent.ChildCount++
			subjectIdNameMap[int64(subject.ParentId)] = parent
		}
	}

	var subjectIdName []vo.Subject
	for _, subject := range subjects {
		name := subject.Name
		if subject.ParentId != 0 {
			name = subjectIdNameMap[int64(subject.ParentId)].Name + name
		}
		subjectIdName = append(subjectIdName, vo.Subject{
			ID:         subject.Id,
			Name:       name,
			ParentId:   subject.ParentId,
			ChildCount: subjectIdNameMap[subject.Id].ChildCount,
		})

	}

	return subjectIdName, nil
}
