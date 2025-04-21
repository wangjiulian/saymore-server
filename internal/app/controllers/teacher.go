package controllers

import (
	"com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/models/dto"
	"com.say.more.server/internal/app/models/vo"
	"com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/types"
	"com.say.more.server/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func TeacherRecommends(c *gin.Context) {

	var teachers []db.Teacher
	if err := repository.Repos.DB.Where(map[string]string{
		db.TeacherFields.IsActive:    "1",
		db.TeacherFields.IsRecommend: "1",
	}).Select([]string{
		db.TeacherFields.ID,
		db.TeacherFields.Nickname,
		db.TeacherFields.AvatarUrl,
		db.TeacherFields.TeachingStartDate,
		db.TeacherFields.Gender,
		db.TeacherFields.EducationSchool,
		db.TeacherFields.EducationLevel,
		db.TeacherFields.CourseHours,
	}).Find(&teachers).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// Get teacher IDs
	var teacherIds []int64
	for _, teacher := range teachers {
		teacherIds = append(teacherIds, teacher.Id)
	}

	teacherIdSubjectMap, err := getTeacherSubjectsMap(teacherIds)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	var recommends []vo.Recommends
	for _, teacher := range teachers {
		var teachingStartDate = "Unknown"
		if teacher.TeachingStartDate.Format("2006-01-02") != "0001-01-01" {
			years := utils.YearsBetween(time.Now(), teacher.TeachingStartDate)
			teachingStartDate = fmt.Sprintf("%d years", years)
		}

		recommends = append(recommends, vo.Recommends{
			ID:                teacher.Base.Id,
			Name:              teacher.Name,
			Nickname:          teacher.Nickname,
			AvatarUrl:         teacher.AvatarUrl,
			TeachingStartDate: teachingStartDate,
			Gender:            db.GenderMap[teacher.Gender],
			EducationSchool:   teacher.EducationSchool,
			EducationLevel:    db.EducationLevelMap[teacher.EducationLevel],
			CourseNum:         teacher.CourseHours.String(),
			SubjectIDs:        teacherIdSubjectMap[teacher.Id],
		})
	}

	SuccessWithData(c, recommends)
}

func TeacherSearch(c *gin.Context) {

	search := strings.TrimSpace(c.Query("search"))
	gender := strings.TrimSpace(c.Query("gender"))
	startDate := strings.TrimSpace(c.Query("start_date"))
	endDate := strings.TrimSpace(c.Query("end_date"))
	subjectIds := strings.TrimSpace(c.Query("subject_ids"))

	ctx := types.From(c)
	offset := (ctx.Page() - 1) * ctx.PageSize()

	// TODO: Validate parameter format

	query := repository.Repos.DB.Table(db.Teacher{}.TableName()).Where(db.TeacherFields.IsActive, "1")
	if search != "" {
		query = query.Where(db.TeacherFields.Nickname+" like ?", "%"+search+"%")
	}
	if gender != "" && gender != "0" {
		query = query.Where(db.TeacherFields.Gender, gender)
	}
	if len(subjectIds) > 0 {
		query = query.Joins("left join "+db.TeacherSubject{}.TableName()+" on "+db.TeacherFields.ID+" = "+db.TeacherSubjectFields.TeacherId).Where(db.TeacherSubjectFields.SubjectId+" in (?)", strings.Split(subjectIds, ","))
	}

	if startDate != "" && endDate == "" {
		query = query.Joins("left join "+db.TeacherAvailability{}.TableName()+" on "+db.TeacherFields.ID+" = "+db.TeacherAvailabilityFields.TeacherId).
			Where(db.TeacherAvailabilityFields.StartTime+" >= ?", startDate).
			Where(db.TeacherAvailabilityFields.CourseId+" = ?", 0)
	}
	if startDate == "" && endDate != "" {
		query = query.Joins("left join "+db.TeacherAvailability{}.TableName()+" on "+db.TeacherFields.ID+" = "+db.TeacherAvailabilityFields.TeacherId).
			Where(db.TeacherAvailabilityFields.StartTime+" >= ?", fmt.Sprintf("%d", time.Now().Unix())).
			Where(db.TeacherAvailabilityFields.StartTime+" <= ?", endDate).
			Where(db.TeacherAvailabilityFields.CourseId+" = ?", 0)
	}
	if startDate != "" && endDate != "" {
		query = query.Joins("left join "+db.TeacherAvailability{}.TableName()+" on "+db.TeacherFields.ID+" = "+db.TeacherAvailabilityFields.TeacherId).
			Where(db.TeacherAvailabilityFields.StartTime+" >= ?", startDate).
			Where(db.TeacherAvailabilityFields.StartTime+" <= ?", endDate).
			Where(db.TeacherAvailabilityFields.CourseId+" = ?", 0)
	}

	var searchTeachers []dto.SearchTeacher
	if err := query.Distinct(db.TeacherFields.ID).Select([]string{
		db.TeacherFields.ID,
		db.TeacherFields.Nickname,
		db.TeacherFields.AvatarUrl,
		db.TeacherFields.TeachingStartDate,
		db.TeacherFields.Gender,
		db.TeacherFields.EducationSchool,
		db.TeacherFields.EducationLevel,
	}).Select([]string{
		db.TeacherFields.ID,
		db.TeacherFields.Name,
		db.TeacherFields.Nickname,
		db.TeacherFields.AvatarUrl,
		db.TeacherFields.TeachingStartDate,
		db.TeacherFields.Gender,
		db.TeacherFields.EducationSchool,
		db.TeacherFields.EducationLevel,
		db.TeacherFields.CourseHours,
	}).Offset(offset).Limit(ctx.PageSize()).Find(&searchTeachers).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// Get teacher IDs
	var teacherIds []int64
	for _, teacher := range searchTeachers {
		teacherIds = append(teacherIds, teacher.ID)
	}

	teacherIdSubjectMap, err := getTeacherSubjectsMap(teacherIds)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	var list []vo.Recommends
	for _, teacher := range searchTeachers {
		var teachingStartDate = "Unknown"
		if teacher.TeachingStartDate.Format("2006-01-02") != "0001-01-01" {
			years := utils.YearsBetween(time.Now(), teacher.TeachingStartDate)
			teachingStartDate = fmt.Sprintf("%d years", years)
		}

		list = append(list, vo.Recommends{
			ID:                teacher.ID,
			Name:              teacher.Name,
			Nickname:          teacher.Nickname,
			AvatarUrl:         teacher.AvatarUrl,
			TeachingStartDate: teachingStartDate,
			Gender:            db.GenderMap[teacher.Gender],
			EducationSchool:   teacher.EducationSchool,
			EducationLevel:    db.EducationLevelMap[teacher.EducationLevel],
			CourseNum:         teacher.CourseHours.String(),
			SubjectIDs:        teacherIdSubjectMap[teacher.ID],
		})
	}

	result := vo.ListData{
		PageInfo: vo.PageInfo{
			Page:     ctx.Page(),
			PageSize: ctx.PageSize(),
		},
		List: list,
	}

	SuccessWithData(c, result)
}

func TeacherDetail(c *gin.Context) {

	teacherIdStr := c.Param("id")
	// Convert string to int
	teacherId, err := strconv.ParseInt(teacherIdStr, 10, 64)
	if err != nil {
		BadRequest(c, "Invalid ID")
		return
	}
	if teacherId < 1 {
		BadRequest(c, "Id cannot be less than 1")
		return
	}

	var teacher db.Teacher
	if err := repository.Repos.DB.Where(db.TeacherFields.ID, teacherId).
		Omit(db.TeacherFields.CreatedAt, db.TeacherFields.UpdatedAt).
		First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "Teacher not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	var subjects []db.TeacherSubject
	if err := repository.Repos.DB.Table(db.TeacherSubject{}.TableName()).Where(db.TeacherSubjectFields.TeacherId, teacherId).Find(&subjects).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}
	var subjectIds []string
	for _, subject := range subjects {
		subjectIds = append(subjectIds, strconv.Itoa(subject.SubjectId))
	}

	var teachingStartDate = "Unknown"
	if teacher.TeachingStartDate.Format("2006-01-02") != "0001-01-01" {
		years := utils.YearsBetween(time.Now(), teacher.TeachingStartDate)
		teachingStartDate = fmt.Sprintf("%d years", years)
	}
	result := vo.TeacherDetail{
		ID:                   teacher.Base.Id,
		Name:                 teacher.Name,
		Nickname:             teacher.Nickname,
		AvatarUrl:            teacher.AvatarUrl,
		TeachingYears:        teachingStartDate,
		Background:           teacher.Background,
		Gender:               db.GenderMap[teacher.Gender],
		EducationSchool:      teacher.EducationSchool,
		EducationLevel:       db.EducationLevelMap[teacher.EducationLevel],
		CourseNum:            teacher.CourseHours.String(),
		TeachingExperience:   teacher.TeachingExperience,
		TeachingAchievements: teacher.TeachingAchievements,
		Notes:                teacher.Notes,
		SuccessCases:         teacher.SuccessCases,
		SubjectIDs:           subjectIds,
		Evaluation:           teacher.Evaluation.String(),
	}

	SuccessWithData(c, result)
}

func TeacherAvailabilities(c *gin.Context) {
	startDate := strings.TrimSpace(c.Query("start_date"))
	endDate := strings.TrimSpace(c.Query("end_date"))
	teacherIdStr := c.Param("id")
	// Convert string to int
	teacherId, err := strconv.ParseInt(teacherIdStr, 10, 64)
	if err != nil {
		BadRequest(c, "Invalid ID")
		return
	}
	if teacherId < 1 {
		BadRequest(c, "Id cannot be less than 1")
		return
	}

	var availabilities []db.TeacherAvailability

	query := repository.Repos.DB.Where(db.TeacherAvailabilityFields.TeacherId, teacherId).
		Where(db.TeacherAvailabilityFields.CourseId, 0)
	if startDate != "" && endDate == "" {
		query = query.Where(db.TeacherAvailabilityFields.StartTime+" >= ?", startDate)
	}
	if startDate == "" && endDate != "" {
		query = query.Where(db.TeacherAvailabilityFields.StartTime+" <= ?", endDate).
			Where(db.TeacherAvailabilityFields.StartTime+" >= ?", fmt.Sprintf("%d", time.Now().Unix()))
	}
	if startDate != "" && endDate != "" {
		query = query.Where(db.TeacherAvailabilityFields.StartTime+" <= ?", endDate).
			Where(db.TeacherAvailabilityFields.StartTime+" >= ?", startDate)
	}
	if startDate == "" && endDate == "" {
		query = query.Where(db.TeacherAvailabilityFields.StartTime+" >= ?", fmt.Sprintf("%d", time.Now().Unix()))
	}

	if err := query.Find(&availabilities).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	var result []vo.TeacherAvailability
	for _, availability := range availabilities {
		result = append(result, vo.TeacherAvailability{
			ID:        availability.Id,
			StartTime: availability.StartTime,
			EndTime:   availability.EndTime,
		})
	}

	SuccessWithData(c, result)
}

func getTeacherSubjectsMap(teacherIds []int64) (map[int64][]int, error) {
	var teacherIdSubjectMap = make(map[int64][]int)
	if len(teacherIds) < 1 {
		return teacherIdSubjectMap, nil
	}
	var teacherSubjects []db.TeacherSubject
	if err := repository.Repos.DB.Where(db.TeacherSubjectFields.TeacherId+" in (?)", teacherIds).Find(&teacherSubjects).Error; err != nil {
		return nil, err
	}
	if len(teacherSubjects) > 0 {
		for _, teacherSubject := range teacherSubjects {
			teacherIdSubjectMap[teacherSubject.TeacherId] = append(teacherIdSubjectMap[teacherSubject.TeacherId], teacherSubject.SubjectId)
		}
	}
	return teacherIdSubjectMap, nil
}
