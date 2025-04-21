package job

import (
	db2 "com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type StartCourseJob struct {
	log *logrus.Entry
}

func NewStartCourseJob() *StartCourseJob {
	return &StartCourseJob{
		log: repository.Repos.Logger.WithField("job", "StartCourseJob"),
	}
}

func (j StartCourseJob) Name() string {
	return "NewStartCourseJob"
}

func (j StartCourseJob) Cron() string {
	return repository.Repos.Config.Crons.CronStartCourse

}

func (j StartCourseJob) Run() {
	if !repository.Repos.Config.AliTextMsg.Enable {
		return
	}

	j.log.Debug("start start course")

	currentTime := time.Now().Unix()
	// Query courses that have not been pushed yet
	db := repository.Repos.DB
	var notifications []db2.CourseNotification

	noticeTIme := currentTime + (repository.Repos.Config.Crons.CronScanCourseBefore * 60)
	if err := db.Model(db2.CourseNotification{}).Where(map[string]interface{}{
		db2.CourseNotificationTableFields.Status: db2.NoticeStatusNotNotice,
	}).Where(db2.CourseNotificationTableFields.StartTime+" <= ?", noticeTIme).Limit(1000).Find(&notifications).Error; err != nil {
		j.log.WithError(err).Error("query course notification error")
		return
	}

	if len(notifications) < 1 {
		j.log.Debug("no course notification")
		return
	}

	var courseIds []int64
	for _, notification := range notifications {
		courseIds = append(courseIds, notification.CourseId)
	}

	var courses []db2.Course
	if err := db.Model(db2.Course{}).Select([]string{
		db2.CourseFields.ID,
		db2.CourseFields.Status,
	}).Where(db2.CourseFields.ID+" IN ?", courseIds).Find(&courses).Error; err != nil {
		j.log.WithError(err).Error("query course error")
		return
	}

	availableCourseIdMap := make(map[int64]bool)
	var notAvailableCourseIds []int64
	for _, course := range courses {
		if course.Status == db2.CourseStatusWaiting {
			availableCourseIdMap[course.Id] = true
		} else {
			notAvailableCourseIds = append(notAvailableCourseIds, course.Id)
		}
	}

	if len(notAvailableCourseIds) > 0 {
		// Ignore courses that do not need to be pushed
		if err := db.Model(db2.CourseNotification{}).Where(db2.CourseNotificationTableFields.CourseId+" IN ?", notAvailableCourseIds).Updates(map[string]interface{}{
			db2.CourseNotificationTableFields.Status:    db2.NoticeStatusIgnoreNotice,
			db2.CourseNotificationTableFields.UpdatedAt: currentTime,
		}).Error; err != nil {
			j.log.WithError(err).Error("delete course notification error")
		}
	}

	if len(availableCourseIdMap) < 1 {
		j.log.Debug("no available course notification")
		return
	}

	var studentIds, teacherIds []int64
	// TODO: Consider multi-processing for large volumes
	// Query teacher and student IDs
	for _, notification := range notifications {
		if _, ok := availableCourseIdMap[notification.CourseId]; ok {
			teacherIds = append(teacherIds, notification.TeacherId)
			studentIds = append(studentIds, notification.StudentId)
		}
	}

	var students []db2.Student
	if err := db.Model(db2.Student{}).Select([]string{
		db2.StudentFields.ID,
		db2.StudentFields.Nickname,
		db2.StudentFields.Phone,
	}).Where(db2.StudentFields.ID+" IN ?", studentIds).Find(&students).Error; err != nil {
		j.log.WithError(err).Error("query student error")
		return
	}
	var studentMap = make(map[int64]db2.Student)
	for _, student := range students {
		studentMap[student.Id] = student
	}

	var teachers []db2.Teacher
	if err := db.Model(db2.Teacher{}).Select([]string{
		db2.TeacherFields.ID,
		db2.TeacherFields.Name,
		db2.TeacherFields.Phone,
	}).Where(db2.TeacherFields.ID+" IN ?", teacherIds).Find(&teachers).Error; err != nil {
		j.log.WithError(err).Error("query teacher error")
		return
	}
	var teacherMap = make(map[int64]db2.Teacher)
	for _, teacher := range teachers {
		teacherMap[teacher.Id] = teacher
	}

	for _, notification := range notifications {
		if _, ok := availableCourseIdMap[notification.CourseId]; ok {
			// TODO: Notify teacher
			j.log.Infof("Notify teacher: %s phone: %s", teacherMap[notification.TeacherId].Name, studentMap[notification.StudentId].Phone)

			repository.Repos.AliTextMsg.SendSms(studentMap[notification.StudentId].Phone, "123459", "Avata", "SMS_237065484")
			// TODO: Notify student
			j.log.Infof("Notify student: %s phone: %s", studentMap[notification.StudentId].Nickname, studentMap[notification.StudentId].Phone)
			// TODO: Update notification table

			if err := db.Model(db2.CourseNotification{}).Where(db2.CourseNotificationTableFields.CourseId, notification.CourseId).Updates(map[string]interface{}{
				db2.CourseNotificationTableFields.Status:    db2.NoticeStatusNotice,
				db2.CourseNotificationTableFields.UpdatedAt: currentTime,
			}).Error; err != nil {
				j.log.WithError(err).Error("update course notification error")
			}
		}
	}

	j.log.Debug("end start course notification")

}
