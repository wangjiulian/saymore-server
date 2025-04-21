package db

const CourseOperationLogTableName = "course_operation_logs"

// 操作员类型：1=教务老师，2=老师，3=学生
const (
	OperatorTypeAdminTeacher = iota + 1
	OperatorTypeTeacher
	OperatorTypeStudent
)

// CourseOperationLog 课程记录操作日志表
type CourseOperationLog struct {
	Base
	CourseId      int64  `gorm:"column:course_id;type:bigint(20);default:0;comment:关联的课程记录表ID;NOT NULL" json:"course_id"`
	OperatorId    int64  `gorm:"column:operator_id;type:bigint(20);default:0;comment:操作人ID;NOT NULL" json:"operator_id"`
	OperatorType  int    `gorm:"column:operator_type;type:tinyint(1);default:0;comment:操作员类型：1=教务老师，2=老师，3=学生;NOT NULL" json:"operator_type"`
	OldStatus     int    `gorm:"column:old_status;type:tinyint(1);default:0;comment:操作前的课程状态;NOT NULL" json:"old_status"`
	NewStatus     int    `gorm:"column:new_status;type:tinyint(1);default:0;comment:操作后的课程状态;NOT NULL" json:"new_status"`
	OperationTime int    `gorm:"column:operation_time;type:timestamp;default:0;comment:操作时间戳;NOT NULL" json:"operation_time"`
	CancelReason  string `gorm:"column:cancel_reason;type:varchar(255);comment:取消原因;NOT NULL" json:"cancel_reason"`
}

func (CourseOperationLog) TableName() string {
	return CourseOperationLogTableName
}

var CourseOperationLogFields = struct {
	ID            string
	CourseId      string
	OperatorId    string
	OperatorType  string
	OldStatus     string
	NewStatus     string
	OperationTime string
	CancelReason  string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            CourseOperationLogTableName + ".id",
	CourseId:      CourseOperationLogTableName + ".course_id",
	OperatorId:    CourseOperationLogTableName + ".operator_id",
	OperatorType:  CourseOperationLogTableName + ".operator_type",
	OldStatus:     CourseOperationLogTableName + ".old_status",
	NewStatus:     CourseOperationLogTableName + ".new_status",
	OperationTime: CourseOperationLogTableName + ".operation_time",
	CancelReason:  CourseOperationLogTableName + ".cancel_reason",
	CreatedAt:     CourseOperationLogTableName + ".created_at",
	UpdatedAt:     CourseOperationLogTableName + ".updated_at",
}
