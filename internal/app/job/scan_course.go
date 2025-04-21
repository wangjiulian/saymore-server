package job

import (
	db2 "com.say.more.server/internal/app/models/db"
	"com.say.more.server/internal/app/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type ScanCourseJob struct {
	log *logrus.Entry
}

func NewScanCourseJob() *ScanCourseJob {
	return &ScanCourseJob{
		log: repository.Repos.Logger.WithField("job", "ScanCourseJob"),
	}
}

func (j ScanCourseJob) Name() string {
	return "ScanCourseJob"
}

func (j ScanCourseJob) Cron() string {
	return repository.Repos.Config.Crons.CronScanCourse

}

func (j ScanCourseJob) Run() {
	if !repository.Repos.Config.AliTextMsg.Enable {
		return
	}
	j.log.Debug("start scan course")
	currentTime := time.Now().Unix()
	// Query courses in waiting status
	var courses []db2.Course
	db := repository.Repos.DB

	if err := db.Model(db2.Course{}).Select([]string{
		db2.CourseFields.ID,
		db2.CourseFields.StudentId,
		db2.CourseFields.TeacherId,
		db2.CourseFields.Name,
		db2.CourseFields.StartTime,
		db2.CourseFields.EndTime,
	}).Where(db2.CourseFields.Status, db2.CourseStatusWaiting).Find(&courses).Limit(1000).Error; err != nil {
		j.log.WithError(err).Error("query waiting course error")
		return
	}

	if len(courses) < 1 {
		j.log.Debug("no waiting course")
		return
	}

	var courseIds []int64
	for _, course := range courses {
		courseIds = append(courseIds, course.Id)
	}

	// Query notification table
	var existNotifications []db2.CourseNotification
	if err := db.Model(db2.CourseNotification{}).Select([]string{db2.CourseNotificationTableFields.CourseId}).Where(db2.CourseNotificationTableFields.CourseId+" IN ?", courseIds).Find(&existNotifications).Error; err != nil {
		j.log.WithError(err).Error("query course notification error")
		return
	}

	var existCourseIdMap = make(map[int64]bool)
	for _, notification := range existNotifications {
		existCourseIdMap[notification.CourseId] = true
	}

	// Filter courses that are already in the notification table
	var needNotifications []db2.CourseNotification
	for _, course := range courses {
		if existCourseIdMap[course.Id] {
			continue
		}
		needNotifications = append(needNotifications, db2.CourseNotification{
			CourseId:  course.Id,
			Title:     course.Name,
			StartTime: course.StartTime,
			EndTime:   course.EndTime,
			TeacherId: course.TeacherId,
			StudentId: course.StudentId,
			Base: db2.Base{
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		})

	}

	if len(needNotifications) > 0 {
		if err := db.CreateInBatches(needNotifications, 1000).Error; err != nil {
			j.log.WithError(err).Error("create course notification error")
			return
		}
	}

	j.log.Debug("finish scan course")
}
