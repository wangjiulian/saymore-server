package db

const CourseNotificationTableName = "course_notifications"

const (
	NoticeStatusNotice       = iota + 1 // 已通知
	NoticeStatusNotNotice               // 未通知
	NoticeStatusIgnoreNotice            // 忽略推送
)

// CourseNotification 课程通知表
type CourseNotification struct {
	Base
	TeacherId int64  `gorm:"column:teacher_id;type:bigint(20);default:0;comment:老师ID;NOT NULL" json:"teacher_id"`
	StudentId int64  `gorm:"column:student_id;type:bigint(20);default:0;comment:学生ID;NOT NULL" json:"student_id"`
	CourseId  int64  `gorm:"column:course_id;type:bigint(20);default:0;comment:课程ID;NOT NULL" json:"course_id"`
	Title     string `gorm:"column:title;type:varchar(255);comment:通知标题;NOT NULL" json:"title"`
	StartTime uint64 `gorm:"column:start_time;type:bigint(20) unsigned;default:0;comment:课程开始时间（Unix时间戳，单位：秒）;NOT NULL" json:"start_time"`
	EndTime   uint64 `gorm:"column:end_time;type:bigint(20) unsigned;default:0;comment:课程结束时间（Unix时间戳，单位：秒）;NOT NULL" json:"end_time"`
	Status    int    `gorm:"column:status;type:tinyint(4);default:2;comment:状态 1：已推送 2 未推送 3 忽略推送;NOT NULL" json:"status"`
}

func (CourseNotification) TableName() string {
	return CourseNotificationTableName
}

var CourseNotificationTableFields = struct {
	Id        string
	TeacherId string
	StudentId string
	Title     string
	CourseId  string
	Status    string
	StartTime string
	EndTime   string
	UpdatedAt string
	CreatedAt string
}{
	Id:        CourseNotificationTableName + ".id",
	Title:     CourseNotificationTableName + ".title",
	CourseId:  CourseNotificationTableName + ".course_id",
	TeacherId: CourseNotificationTableName + ".teacher_id",
	StudentId: CourseNotificationTableName + ".student_id",
	StartTime: CourseNotificationTableName + ".start_time",
	EndTime:   CourseNotificationTableName + ".end_time",
	Status:    CourseNotificationTableName + ".status",
	UpdatedAt: CourseNotificationTableName + ".updated_at",
	CreatedAt: CourseNotificationTableName + ".created_at",
}
