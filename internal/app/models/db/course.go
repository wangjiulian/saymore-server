package db

const CourseTableName = "courses"

// 课程类型：1-正式课，2-试听课
const (
	CourseTypeRegular = iota + 1
	CourseTypeTrial
)

// 上课状态：1-待上课，2-上课中，3-已完成，4-已取消
const (
	CourseStatusWaiting = iota + 1
	CourseStatusOnGoing
	CourseStatusFinished
	CourseStatusCanceled
)

// 是否已评价
const (
	CourseIsEvaluated = iota + 1
	CourseIsNotEvaluated
)

var CourseStatusMap = map[int]string{
	CourseStatusWaiting:  "待上课",
	CourseStatusOnGoing:  "上课中",
	CourseStatusFinished: "已完成",
	CourseStatusCanceled: "已取消",
}

var CourseTypeMap = map[int]string{
	CourseTypeRegular: "正式课",
	CourseTypeTrial:   "试听课",
}

// Course 课程记录表
type Course struct {
	Base
	CourseType   int    `gorm:"column:course_type;type:tinyint(1);default:1;comment:课程类型：1-正式课，2-试听课;NOT NULL" json:"course_type"`
	Name         string `gorm:"column:name;type:varchar(100);comment:课程名称;NOT NULL" json:"name"`
	SubjectId    int    `gorm:"column:subject_id;type:int(11);default:0;comment:学科ID;NOT NULL" json:"subject_id"`
	TeacherId    int64  `gorm:"column:teacher_id;type:bigint(20);default:0;comment:老师ID;NOT NULL" json:"teacher_id"`
	StudentId    int64  `gorm:"column:student_id;type:bigint(20);default:0;comment:学生ID;NOT NULL" json:"student_id"`
	Status       int    `gorm:"column:status;type:tinyint(1);default:0;comment:上课状态：1-待上课，2-上课中，3-已完成，4-已取消;NOT NULL" json:"status"`
	IsEvaluated  int    `gorm:"column:is_evaluated;type:tinyint(1);default:2;comment:是否评价: 1-已评价，2-未评价;NOT NULL" json:"status"`
	CancelReason string `gorm:"column:cancel_reason;type:varchar(255);comment:取消原因;NOT NULL" json:"cancel_reason"`
	StartTime    uint64 `gorm:"column:start_time;type:bigint(20) unsigned;default:0;comment:课程开始时间（Unix时间戳，单位：秒）;NOT NULL" json:"start_time"`
	EndTime      uint64 `gorm:"column:end_time;type:bigint(20) unsigned;default:0;comment:课程结束时间（Unix时间戳，单位：秒）;NOT NULL" json:"end_time"`
}

func (Course) TableName() string {
	return CourseTableName
}

var CourseFields = struct {
	ID           string
	CourseType   string
	Name         string
	SubjectId    string
	TeacherId    string
	StudentId    string
	Status       string
	CancelReason string
	StartTime    string
	EndTime      string
	IsEvaluated  string
	CreatedAt    string
	UpdatedAt    string
}{
	ID:           CourseTableName + ".id",
	CourseType:   CourseTableName + ".course_type",
	Name:         CourseTableName + ".name",
	SubjectId:    CourseTableName + ".subject_id",
	TeacherId:    CourseTableName + ".teacher_id",
	StudentId:    CourseTableName + ".student_id",
	Status:       CourseTableName + ".status",
	CancelReason: CourseTableName + ".cancel_reason",
	StartTime:    CourseTableName + ".start_time",
	EndTime:      CourseTableName + ".end_time",
	IsEvaluated:  CourseTableName + ".is_evaluated",
	CreatedAt:    CourseTableName + ".created_at",
	UpdatedAt:    CourseTableName + ".updated_at",
}
