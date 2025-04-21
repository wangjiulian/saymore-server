package controllers

import (
	"com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/models/dto"
	"com.say.more.server/internal/app/models/form"
	"com.say.more.server/internal/app/models/vo"
	"com.say.more.server/internal/app/repository"
	"com.say.more.server/internal/app/types"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// BookRegularCourse books a regular course
func BookRegularCourse(c *gin.Context) {
	currentTime := time.Now().Unix()
	ctx := types.From(c)
	params := &form.BookRegularCourseForm{}
	err := c.ShouldBind(params)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	if params.TeacherId <= 0 || len(params.TeacherAvailabilityIds) <= 0 || params.StudentPackageId <= 0 || params.SubjectId <= 0 {
		BadRequest(c, "invalid params")
		return
	}

	// Parse booking time IDs and deduplicate
	var teacherAvailabilityIds []int64
	var teacherAvailabilityIdMap = make(map[string]struct{})
	for _, id := range strings.Split(params.TeacherAvailabilityIds, ",") {
		id = strings.TrimSpace(id)
		if id == "" {

		}
		if _, ok := teacherAvailabilityIdMap[fmt.Sprintf("%d", id)]; ok {
			continue
		}
		idInt, err := strconv.ParseInt(id, 10, 64)
		if idInt <= 0 {
			BadRequest(c, "invalid teacher availability id")
			return
		}
		if err != nil {
			InternalServerError(c, err.Error())
			return
		}
		teacherAvailabilityIds = append(teacherAvailabilityIds, idInt)
		teacherAvailabilityIdMap[id] = struct{}{}
	}
	if len(teacherAvailabilityIds) <= 0 {
		BadRequest(c, "invalid teacher availability ids")
		return

	}

	tx := repository.Repos.DB

	// Get teacher
	var teacher db.Teacher
	if err := tx.Table(db.Teacher{}.TableName()).Where(db.TeacherFields.ID, params.TeacherId).
		First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher not found")
			return
		}
	}

	// Check if the teacher can teach the subject
	var teacherSubject db.TeacherSubject
	if err := tx.Table(db.TeacherSubject{}.TableName()).Where(db.TeacherSubjectFields.TeacherId, params.TeacherId).
		Where(db.TeacherSubjectFields.SubjectId, params.SubjectId).
		First(&teacherSubject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher subject not found")
			return
		}

		InternalServerError(c, err.Error())
		return
	}

	// Get teacher's available time slots
	var teacherAvailabilities []db.TeacherAvailability
	if err := tx.Table(db.TeacherAvailability{}.TableName()).
		Where(db.TeacherAvailabilityFields.ID, teacherAvailabilityIds).
		Order(db.TeacherAvailabilityFields.StartTime + " asc").Find(&teacherAvailabilities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher availability not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}
	if len(teacherAvailabilities) <= 0 {
		BadRequest(c, "teacher availability not found")
		return

	}
	// Check if teacher's available time slots are booked and if multiple slots are consecutive
	var startTime, endTime, currentAvailabilityTime uint64
	for index, teacherAvailability := range teacherAvailabilities {
		if currentAvailabilityTime == 0 {
			// Course start time
			startTime = teacherAvailability.StartTime
			// Store cursor for consecutive course time
			currentAvailabilityTime = teacherAvailability.EndTime
		} else {
			if currentAvailabilityTime != teacherAvailability.StartTime {
				BadRequest(c, "Teacher's time slots are not consecutive, please exit, refresh, and book again!")
				return
			}

			currentAvailabilityTime = teacherAvailability.EndTime
		}

		if teacherAvailability.CourseId > 0 {
			BadRequest(c, "Teacher's time slot is already booked, please exit, refresh, and book again!")
			return
		}
		if (index + 1) == len(teacherAvailabilities) {
			// Set course end time
			endTime = teacherAvailability.EndTime
		}
	}

	// Get student package
	var studentPackage db.StudentPackage
	if err := tx.Table(db.StudentPackage{}.TableName()).Where(map[string]interface{}{
		db.StudentPackageFields.ID:        params.StudentPackageId,
		db.StudentPackageFields.StudentId: ctx.StudentID(),
	}).First(&studentPackage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "student package not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	// Consume student course hours
	// Number of sessions * hours per session unit
	consumeHours := decimal.NewFromInt(int64(len(teacherAvailabilities))).Mul(decimal.NewFromFloat(repository.Repos.Config.Course.CourseUnit))
	if studentPackage.Hours.LessThan(consumeHours) {
		BadRequest(c, "Insufficient package balance")
		return
	}

	// Query package bookable subject IDs, including current ID
	var subjectList []db.Subject
	if err := tx.Table(db.Subject{}.TableName()).Select(db.SubjectFields.ID, db.SubjectFields.ParentId, db.SubjectFields.Name).Where(db.SubjectFields.ID, studentPackage.SubjectId).Or(db.SubjectFields.ParentId, studentPackage.SubjectId).Find(&subjectList).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}
	if len(subjectList) < 1 {
		BadRequest(c, "package subject not found")
		return
	}

	// Check if current package can book the current subject
	var exist bool
	var childSubject, parentSubject string
	for _, v := range subjectList {
		if v.ParentId == 0 {
			parentSubject = v.Name
		}
		if v.Id == int64(params.SubjectId) {
			childSubject = v.Name
			exist = true
		}
	}

	// Check if teacher can use current package
	if err := tx.Table(db.StudentTeacherPackages{}.TableName()).Where(map[string]interface{}{
		db.StudentTeacherPackageFields.StudentPackageId: params.StudentPackageId,
		db.StudentTeacherPackageFields.TeacherId:        params.TeacherId,
	}).First(&db.StudentTeacherPackages{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher package not found")
			return
		}

		InternalServerError(c, err.Error())
		return
	}

	courseName := parentSubject + childSubject + " course"
	if parentSubject == childSubject {
		courseName = parentSubject + " course"
	}

	if !exist {
		BadRequest(c, "package subject is not available")
		return
	}

	if err := tx.Transaction(func(d *gorm.DB) error {
		// Create course record
		course := db.Course{
			CourseType: db.CourseTypeRegular,
			Name:       courseName,
			SubjectId:  params.SubjectId,
			TeacherId:  teacher.Id,
			StudentId:  ctx.StudentID(),
			Status:     db.CourseStatusWaiting,
			StartTime:  startTime,
			EndTime:    endTime,
			Base: db.Base{
				UpdatedAt: currentTime,
				CreatedAt: currentTime,
			},
		}
		if err := d.Create(&course).Error; err != nil {
			return err
		}

		// Create course operation log
		courseOperationLog := db.CourseOperationLog{
			CourseId:      course.Id,
			OperatorId:    ctx.StudentID(),
			OperatorType:  db.OperatorTypeStudent,
			OldStatus:     0,
			NewStatus:     db.CourseStatusWaiting,
			OperationTime: int(time.Now().Unix()),
			Base: db.Base{
				UpdatedAt: currentTime,
				CreatedAt: currentTime,
			},
		}
		if err := d.Create(&courseOperationLog).Error; err != nil {
			return err
		}
		updatePackage := d.Table(db.StudentPackage{}.TableName()).Where(map[string]interface{}{
			db.StudentPackageFields.ID:        params.StudentPackageId,
			db.StudentPackageFields.StudentId: ctx.StudentID(),
			db.StudentPackageFields.LeftHours: studentPackage.LeftHours,
		}).Updates(map[string]interface{}{
			db.StudentPackageFields.LeftHours: gorm.Expr(db.StudentPackageFields.LeftHours+" - ?", consumeHours.String()),
			db.StudentPackageFields.UpdatedAt: currentTime,
		})
		if updatePackage.Error != nil {
			return updatePackage.Error
		}
		if updatePackage.RowsAffected < 1 {
			return errors.New("update student package failed")
		}

		// Add student package consumption record
		studentPackageDetail := db.StudentPackageDetail{
			StudentId:        ctx.StudentID(),
			StudentPackageId: studentPackage.Id,
			CourseId:         course.Id,
			Title:            "Book " + courseName,
			Hours:            consumeHours,
			LeftHours:        studentPackage.LeftHours.Sub(consumeHours),
			Change:           db.ChangeTypeReduce,
			ChangeType:       db.ChangeTypeBookCourse,
			Base: db.Base{
				UpdatedAt: currentTime,
				CreatedAt: currentTime,
			},
		}
		if err := d.Create(&studentPackageDetail).Error; err != nil {
			return err
		}

		updateTeacherAvailability := d.Table(db.TeacherAvailability{}.TableName()).Where(map[string]interface{}{
			db.TeacherAvailabilityFields.CourseId: 0,
		}).Where(db.TeacherAvailabilityFields.ID+" in(?)", teacherAvailabilityIds).
			Updates(map[string]interface{}{
				db.TeacherAvailabilityFields.CourseId:  course.Id,
				db.TeacherAvailabilityFields.UpdatedAt: currentTime,
			})
		if updateTeacherAvailability.Error != nil {
			return updateTeacherAvailability.Error
		}
		if updateTeacherAvailability.RowsAffected < 1 {
			return errors.New("update teacher availability failed")
		}

		return nil
	}); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c)
}

// BookTrialCourse books a trial course
func BookTrialCourse(c *gin.Context) {
	currentTime := time.Now().Unix()
	ctx := types.From(c)
	params := &form.BookTrialCourseForm{}
	err := c.ShouldBind(params)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	if params.TeacherId <= 0 || params.TeacherAvailabilityId <= 0 || params.SubjectId <= 0 {
		BadRequest(c, "invalid params")
		return
	}

	tx := repository.Repos.DB

	// Get student trial course quota
	var studentTrial db.StudentTrialQuota
	if err := tx.Table(db.StudentTrialQuota{}.TableName()).Where(map[string]interface{}{
		db.StudentTrialQuotaFields.StudentId: ctx.StudentID(),
		db.StudentTrialQuotaFields.CourseId:  0,
	}).First(&studentTrial).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "You don't have enough trial course quota")
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	// Get teacher
	var teacher db.Teacher
	if err := tx.Table(db.Teacher{}.TableName()).Where(db.TeacherFields.ID, params.TeacherId).
		First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher not found")
			return
		}
	}

	// Check if teacher can teach the subject
	var teacherSubject db.TeacherSubject
	if err := tx.Table(db.TeacherSubject{}.TableName()).Where(db.TeacherSubjectFields.TeacherId, params.TeacherId).
		Where(db.TeacherSubjectFields.SubjectId, params.SubjectId).
		First(&teacherSubject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher subject not found")
			return
		}

		InternalServerError(c, err.Error())
		return
	}

	// Get teacher's available time slot
	var teacherAvailability db.TeacherAvailability
	if err := tx.Table(db.TeacherAvailability{}.TableName()).Where(map[string]interface{}{
		db.TeacherAvailabilityFields.ID:       params.TeacherAvailabilityId,
		db.TeacherAvailabilityFields.CourseId: 0,
	}).First(&teacherAvailability).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "teacher availability not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	// Get course name
	var subject db.Subject
	if err := tx.Table(db.Subject{}.TableName()).Where(db.SubjectFields.ID, params.SubjectId).
		First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "subject not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}
	courseName := subject.Name
	if subject.ParentId > 0 {
		// Get parent category name
		var parentSubject db.Subject
		if err := tx.Table(db.Subject{}.TableName()).Where(db.SubjectFields.ID, subject.ParentId).
			First(&parentSubject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				BadRequest(c, "parent subject not found")
				return
			}
			InternalServerError(c, err.Error())
			return
		}
		courseName = parentSubject.Name + courseName + " course"
	}

	if err := tx.Transaction(func(d *gorm.DB) error {
		// Create course record
		course := db.Course{
			CourseType: db.CourseTypeTrial,
			Name:       courseName,
			SubjectId:  params.SubjectId,
			TeacherId:  teacher.Id,
			StudentId:  ctx.StudentID(),
			Status:     db.CourseStatusWaiting,
			StartTime:  teacherAvailability.StartTime,
			EndTime:    teacherAvailability.EndTime,
			Base: db.Base{
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		}
		if err := d.Create(&course).Error; err != nil {
			return err
		}

		// Create course operation log
		courseOperationLog := db.CourseOperationLog{
			CourseId:      course.Id,
			OperatorId:    ctx.StudentID(),
			OperatorType:  db.OperatorTypeStudent,
			OldStatus:     0,
			NewStatus:     db.CourseStatusWaiting,
			OperationTime: int(time.Now().Unix()),
			Base: db.Base{
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		}
		if err := d.Create(&courseOperationLog).Error; err != nil {
			return err
		}

		// Consume student trial quota
		updateStudentTrialQuota := d.Table(db.StudentTrialQuota{}.TableName()).Where(map[string]interface{}{
			db.StudentTrialQuotaFields.StudentId: ctx.StudentID(),
			db.StudentTrialQuotaFields.CourseId:  0,
		}).Update(db.StudentTrialQuotaFields.CourseId, course.Id)
		if updateStudentTrialQuota.Error != nil {
			return updateStudentTrialQuota.Error
		}
		if updateStudentTrialQuota.RowsAffected < 1 {
			return errors.New("book regular course failed")
		}

		updateTeacherAvailability := d.Table(db.TeacherAvailability{}.TableName()).Where(map[string]interface{}{
			db.TeacherAvailabilityFields.ID:       teacherAvailability.Id,
			db.TeacherAvailabilityFields.CourseId: 0,
		}).Update(db.TeacherAvailabilityFields.CourseId, course.Id)
		if updateTeacherAvailability.Error != nil {
			return updateTeacherAvailability.Error
		}
		if updateTeacherAvailability.RowsAffected < 1 {
			return errors.New("book regular course failed")
		}

		return nil
	}); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c)
}

// CancelCourse cancels a course
func CancelCourse(c *gin.Context) {
	current := time.Now().Unix()
	ctx := types.From(c)
	courseIdStr := c.Param("course_id")
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil || courseId < 1 {
		BadRequest(c, "invalid course id")
		return
	}
	params := &form.CancelCourseForm{}
	if err := c.ShouldBind(params); err != nil {
		BadRequest(c, err.Error())
		return
	}
	params.Reason = strings.TrimSpace(params.Reason)
	if params.Reason == "" {
		BadRequest(c, "invalid reason")
		return
	}

	tx := repository.Repos.DB

	// Query course
	var course db.Course
	if err := tx.Where(map[string]interface{}{
		db.CourseFields.ID:        courseId,
		db.CourseFields.StudentId: ctx.StudentID(),
	}).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BadRequest(c, "course not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}
	if course.Status == db.CourseStatusCanceled || course.Status == db.CourseStatusFinished {
		BadRequest(c, "course already canceled or finished")
		return
	}

	if err := tx.Transaction(func(d *gorm.DB) error {
		// Update course record status
		updateCourse := d.Table(db.Course{}.TableName()).Where(map[string]interface{}{
			db.CourseFields.ID:        course.Id,
			db.CourseFields.StudentId: ctx.StudentID(),
			db.CourseFields.Status:    course.Status,
		}).Update(db.CourseFields.Status, db.CourseStatusCanceled)
		if updateCourse.Error != nil {
			return updateCourse.Error
		}
		if updateCourse.RowsAffected < 1 {
			return errors.New("cancel course failed")
		}
		// Create course operation log
		courseOperationLog := db.CourseOperationLog{
			CourseId:      course.Id,
			OperatorId:    ctx.StudentID(),
			OperatorType:  db.OperatorTypeStudent,
			OldStatus:     course.Status,
			NewStatus:     db.CourseStatusCanceled,
			OperationTime: int(time.Now().Unix()),
			CancelReason:  params.Reason,
			Base: db.Base{
				CreatedAt: current,
				UpdatedAt: current,
			},
		}
		if err := d.Create(&courseOperationLog).Error; err != nil {
			return err
		}

		// Update teacher availability record
		updateTeacherAvailability := d.Table(db.TeacherAvailability{}.TableName()).Where(map[string]interface{}{
			db.TeacherAvailabilityFields.TeacherId: course.TeacherId,
			db.TeacherAvailabilityFields.CourseId:  course.Id,
		}).Update(db.TeacherAvailabilityFields.CourseId, 0)
		if updateTeacherAvailability.Error != nil {
			return updateTeacherAvailability.Error
		}
		if updateTeacherAvailability.RowsAffected < 1 {
			return errors.New("update teacher availability failed")
		}

		// Check course type
		switch course.CourseType {
		case db.CourseTypeRegular:
			// For regular course type, consider refund based on cancellation responsibility
			cancelInterval := repository.Repos.Config.Course.CancelInterval * 60
			ChangeTypeCancel := db.ChangeTypeCancelCourseResponsible // Default: responsible cancellation
			cancelRefund := decimal.NewFromFloat(0.5)                // Default: refund 0.5 hours
			// Check if current cancellation time is beyond the configuration time interval from course start time
			if int64(course.StartTime)-current >= int64(cancelInterval) {
				// Within no-responsibility cancellation range
				cancelRefund = decimal.NewFromFloat(1) // No-responsibility cancellation, full refund
				ChangeTypeCancel = db.ChangeTypeCancelCourseNoResponsible
			}

			// Query course package used
			var studentPackageDetail db.StudentPackageDetail
			if err := d.Where(map[string]interface{}{
				db.StudentPackageDetailFields.StudentId: ctx.StudentID(),
				db.StudentPackageDetailFields.CourseId:  course.Id,
			}).First(&studentPackageDetail).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("invalid course package")
				}
				return err
			}

			// Query current package
			var studentPackage db.StudentPackage
			if err := d.Where(map[string]interface{}{
				db.StudentPackageFields.ID:        studentPackageDetail.StudentPackageId,
				db.StudentPackageFields.StudentId: ctx.StudentID(),
			}).First(&studentPackage).Error; err != nil {
				return err
			}

			totalRefund := cancelRefund
			// Update student package
			updateStudentPackage := d.Table(db.StudentPackage{}.TableName()).Where(map[string]interface{}{
				db.StudentPackageFields.ID:        studentPackageDetail.StudentPackageId,
				db.StudentPackageFields.StudentId: ctx.StudentID(),
				db.StudentPackageFields.LeftHours: studentPackage.LeftHours,
			}).Update(db.StudentPackageFields.LeftHours, gorm.Expr(db.StudentPackageFields.LeftHours+" + ?", totalRefund))
			if updateStudentPackage.Error != nil {
				return updateStudentPackage.Error
			}
			if updateStudentPackage.RowsAffected < 1 {
				return errors.New("Failed to cancel course, please try again later")
			}

			// Add student package consumption record
			studentPackageDetailAdd := db.StudentPackageDetail{
				StudentId:        ctx.StudentID(),
				StudentPackageId: studentPackageDetail.StudentPackageId,
				CourseId:         course.Id,
				Title:            "Cancel " + course.Name,
				Hours:            totalRefund,
				LeftHours:        studentPackage.LeftHours.Add(totalRefund),
				Change:           db.ChangeTypeAdd,
				ChangeType:       uint(ChangeTypeCancel),
				Base: db.Base{
					CreatedAt: current,
					UpdatedAt: current,
				},
			}
			if err := d.Create(&studentPackageDetailAdd).Error; err != nil {
				return err
			}

		case db.CourseTypeTrial:
			// For trial course type, only need to refund trial opportunity
			updateStudentTrialQuota := d.Table(db.StudentTrialQuota{}.TableName()).Where(map[string]interface{}{
				db.StudentTrialQuotaFields.StudentId: ctx.StudentID(),
				db.StudentTrialQuotaFields.CourseId:  course.Id,
			}).Update(db.StudentTrialQuotaFields.CourseId, 0)
			if updateStudentTrialQuota.Error != nil {
				return updateStudentTrialQuota.Error
			}
			if updateStudentTrialQuota.RowsAffected < 1 {
				return errors.New("update student trial quota failed")
			}

		default:
			return errors.New("invalid course type")
		}

		return nil
	}); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c)
}

// Courses gets the course list
func Courses(c *gin.Context) {
	ctx := types.From(c)
	offset := (ctx.Page() - 1) * ctx.PageSize()
	statusStr := c.DefaultQuery("status", "0")
	status, err := strconv.Atoi(statusStr)
	if err != nil || status < 0 {
		BadRequest(c, "invalid status")
		return
	}
	if status > 0 {
		if _, ok := db.CourseStatusMap[status]; !ok {
			BadRequest(c, "invalid status")
			return
		}
	}

	var courses []db.Course
	query := repository.Repos.DB.Where(map[string]interface{}{
		db.CourseFields.StudentId: ctx.StudentID(),
	})
	if status > 0 {
		query = query.Where(db.CourseFields.Status, status)
	}
	if err := query.Offset(offset).Limit(ctx.PageSize()).Order(db.CourseFields.ID + " desc").Find(&courses).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// Get teacher names
	var teacherMap = make(map[int64]db.Teacher)
	var teacherIds []int64
	for _, course := range courses {
		teacherIds = append(teacherIds, course.TeacherId)
	}
	if teacherIds != nil && len(teacherIds) > 0 {
		var teachers []db.Teacher
		if err := repository.Repos.DB.Where(db.TeacherFields.ID+" in ?", teacherIds).Find(&teachers).Error; err != nil {
			InternalServerError(c, err.Error())
			return
		}
		for _, teacher := range teachers {
			teacherMap[teacher.Id] = teacher
		}
	}

	var list []vo.CourseList
	for _, course := range courses {
		list = append(list, vo.CourseList{
			Id:             course.Id,
			TeacherId:      course.TeacherId,
			SubjectId:      course.SubjectId,
			Name:           course.Name,
			CourseTypeName: db.CourseTypeMap[course.CourseType],
			TeacherName:    teacherMap[course.TeacherId].Nickname,
			Status:         course.Status,
			IsEvaluated:    course.IsEvaluated,
			StartTime:      course.StartTime,
			EndTime:        course.EndTime,
		})

	}

	var result vo.ListData
	result.PageInfo = vo.PageInfo{
		Page:     ctx.Page(),
		PageSize: ctx.PageSize(),
	}
	result.List = list
	SuccessWithData(c, result)
}

func CourseDetail(c *gin.Context) {

	courseIdStr := c.Param("course_id")
	// Convert string to int
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil {
		BadRequest(c, "Invalid ID")
		return
	}
	if courseId < 1 {
		BadRequest(c, "Id cannot be less than 1")
		return
	}

	var course db.Course
	if err := repository.Repos.DB.Where(db.CourseFields.ID, courseId).
		First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "Course not found")
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	// Get teacher details
	var teacher db.Teacher
	if err := repository.Repos.DB.Where(db.TeacherFields.ID, course.TeacherId).First(&teacher).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	courseDetail := vo.CourseDetail{
		Id:                     course.Id,
		SubjectId:              course.SubjectId,
		Name:                   course.Name,
		CourseTypeName:         db.CourseTypeMap[course.CourseType],
		TeacherId:              course.TeacherId,
		TeacherName:            teacher.Name,
		TeacherNickName:        teacher.Nickname,
		TeacherAvatar:          teacher.AvatarUrl,
		TeacherEducationSchool: teacher.EducationSchool,
		TeacherEducationLevel:  db.EducationLevelMap[teacher.EducationLevel],
		Status:                 course.Status,
		StatusText:             db.CourseStatusMap[course.Status],
		StartTime:              course.StartTime,
		EndTime:                course.EndTime,
		IsEvaluated:            course.IsEvaluated,
		CreatedTime:            course.CreatedAt,
	}

	SuccessWithData(c, courseDetail)
}

// CourseEvaluations gets course evaluations
func CourseEvaluations(c *gin.Context) {
	var teacherId int64
	ctx := types.From(c)
	offset := (ctx.Page() - 1) * ctx.PageSize()
	teacherIdStr := c.DefaultQuery("teacher_id", "")
	if teacherIdStr != "" {
		parseTeacherId, err := strconv.ParseInt(teacherIdStr, 10, 64)
		if err != nil {
			BadRequest(c, "Invalid teacher_id")
			return
		}
		teacherId = parseTeacherId
	}

	var courseEvaluations []db.CourseEvaluation
	query := repository.Repos.DB.Table(db.CourseEvaluation{}.TableName())
	if teacherId > 0 {
		query = query.Where(db.CourseEvaluationFields.TeacherId, teacherId)

	}
	if err := query.Offset(offset).Limit(ctx.PageSize()).Order(db.CourseEvaluationFields.CreatedAt + " desc").Find(&courseEvaluations).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// Get student names
	var studentMap = make(map[int64]db.Student)
	var studentIds []int64
	for _, evaluation := range courseEvaluations {
		studentIds = append(studentIds, evaluation.StudentId)
	}
	if studentIds != nil && len(studentIds) > 0 {
		var students []db.Student
		if err := repository.Repos.DB.Where(db.StudentFields.ID+" in ?", studentIds).Find(&students).Error; err != nil {
			InternalServerError(c, err.Error())
			return
		}
		for _, student := range students {
			studentMap[student.Id] = student
		}
	}

	var list []vo.CourseEvaluationList
	for _, evaluation := range courseEvaluations {
		avgRating := decimal.NewFromInt(int64(evaluation.ContentQualityRating + evaluation.InstructorClarityRating + evaluation.LearningGainRating)).Div(decimal.NewFromInt(3))
		avgRating = avgRating.RoundUp(1)
		list = append(list, vo.CourseEvaluationList{
			Id:            evaluation.Id,
			StudentName:   studentMap[evaluation.StudentId].Nickname,
			StudentAvatar: studentMap[evaluation.StudentId].AvatarUrl,
			AvgRating:     avgRating,
			Content:       evaluation.Content,
			CreatedTime:   evaluation.CreatedAt,
		})

	}

	var result vo.ListData
	result.PageInfo = vo.PageInfo{
		Page:     ctx.Page(),
		PageSize: ctx.PageSize(),
	}
	result.List = list
	SuccessWithData(c, result)
}

func AddCourseEvaluation(c *gin.Context) {
	currentTime := time.Now().Unix()
	courseIdStr := c.Param("course_id")
	// Convert string to int
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil {
		BadRequest(c, "Invalid ID")
		return
	}
	if courseId < 1 {
		BadRequest(c, "Id cannot be less than 1")
		return
	}
	ctx := types.From(c)
	params := &form.CourseEvaluationForm{}
	if err := c.ShouldBind(params); err != nil {
		BadRequest(c, err.Error())
		return
	}
	if params.ContentQualityRating <= 0 || params.InstructorClarityRating <= 0 || params.LearningGainRating <= 0 {
		BadRequest(c, "invalid params")
		return
	}

	var course db.Course
	if err := repository.Repos.DB.First(&course, courseId).Error; err != nil {
		NotFound(c, "Course not found")
		return
	}
	if course.StudentId != ctx.StudentID() {
		BadRequest(c, "not your course")
		return
	}

	if course.IsEvaluated == db.CourseIsEvaluated {
		// Already evaluated
		BadRequest(c, "already evaluated")
		return
	}

	if err := repository.Repos.DB.Transaction(func(tx *gorm.DB) error {
		// Create evaluation
		if err := tx.Create(&db.CourseEvaluation{
			CourseId:                courseId,
			TeacherId:               course.TeacherId,
			StudentId:               course.StudentId,
			ContentQualityRating:    params.ContentQualityRating,
			InstructorClarityRating: params.InstructorClarityRating,
			LearningGainRating:      params.LearningGainRating,
			Content:                 params.Content,
			Base: db.Base{
				UpdatedAt: currentTime,
				CreatedAt: currentTime,
			},
		}).Error; err != nil {
			InternalServerError(c, err.Error())
			return err
		}

		// Calculate teacher's average rating
		var avgRating decimal.Decimal
		var courseAvgRatings dto.CourseAvgRatings
		if err := tx.Table(db.CourseEvaluation{}.TableName()).
			Where(db.CourseEvaluationFields.TeacherId, course.TeacherId).
			Select("AVG(content_quality_rating) AS avg_content_quality, AVG(instructor_clarity_rating) AS avg_clarity, AVG(learning_gain_rating) AS avg_learning_gain").
			Scan(&courseAvgRatings).Error; err != nil {
			InternalServerError(c, err.Error())
			return err
		}
		if courseAvgRatings.AvgClarity.IsZero() {
			// Default won't be zero, so if it's zero, it means there's no evaluation. Use current evaluation as average rating.
			avgRating = decimal.NewFromInt(int64(params.ContentQualityRating + params.InstructorClarityRating + params.LearningGainRating)).Div(decimal.NewFromInt(3))
		} else {
			avgRating = courseAvgRatings.AvgContentQuality.Add(courseAvgRatings.AvgClarity).Add(courseAvgRatings.AvgLearningGain).Div(decimal.NewFromInt(3))
		}

		avgRating = avgRating.RoundUp(1)
		// Update teacher's average rating
		if err := tx.Model(&db.Teacher{}).Where(db.TeacherFields.ID, course.TeacherId).Update(db.TeacherFields.Evaluation, avgRating).Error; err != nil {
			InternalServerError(c, err.Error())
			return err
		}

		// Update course status
		if err := tx.Model(&db.Course{}).Where(db.CourseFields.ID, courseId).Where(db.CourseFields.IsEvaluated, db.CourseIsNotEvaluated).Update(db.CourseFields.IsEvaluated, db.CourseIsEvaluated).Error; err != nil {
			InternalServerError(c, err.Error())
			return err
		}

		return nil
	}); err != nil {
		return
	}

	Success(c)

}

func CourseEvaluationDetail(c *gin.Context) {
	courseIdStr := c.Param("course_id")
	// Convert string to int
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil {
		BadRequest(c, "Invalid ID")
		return
	}
	if courseId < 1 {
		BadRequest(c, "Id cannot be less than 1")
		return
	}
	var courseEvaluation db.CourseEvaluation
	if err := repository.Repos.DB.
		Where(map[string]interface{}{db.CourseEvaluationFields.CourseId: courseId, db.CourseEvaluationFields.StudentId: types.From(c).StudentID()}).First(&courseEvaluation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "Course evaluation not found")
			return
		}
		return
	}

	SuccessWithData(c, courseEvaluation)
}
