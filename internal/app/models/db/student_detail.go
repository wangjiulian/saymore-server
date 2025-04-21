package db

import "github.com/shopspring/decimal"

const StudentPackageDetailTableName = "student_package_details"

const (
	ChangeTypeAdd    = 1
	ChangeTypeReduce = 2
)

const (
	ChangeTypePurchase = iota + 1
	ChangeTypeBookCourse
	ChangeTypeCancelCourseNoResponsible
	ChangeTypeCancelCourseResponsible
)

// StudentPackageDetail 学生课时包消费明细表
type StudentPackageDetail struct {
	Base
	StudentId        int64           `gorm:"column:student_id;type:bigint(20);default:0;comment:关联学生表的 ID;NOT NULL" json:"student_id"`
	StudentPackageId int64           `gorm:"column:student_package_id;type:bigint(20);default:0;comment:关联课时包的 ID;NOT NULL" json:"student_package_id"`
	CourseId         int64           `gorm:"column:course_id;type:bigint(20);default:0;comment:关联课程的 ID;NOT NULL" json:"bind_id"`
	Title            string          `gorm:"column:title;type:varchar(255);comment:课时包标题;NOT NULL" json:"title"`
	Hours            decimal.Decimal `gorm:"column:hours;type:decimal(6,2) unsigned;default:0.0;comment:课时数量;NOT NULL" json:"hours"`
	LeftHours        decimal.Decimal `gorm:"column:left_hours;type:decimal(6,2) unsigned;default:0.0;comment:剩余课时包数量;NOT NULL" json:"left_hours"`
	Change           uint            `gorm:"column:change;type:tinyint(4) unsigned;default:0;comment:课时变化 1：增加 2：减少;NOT NULL" json:"change"`
	ChangeType       uint            `gorm:"column:change_type;type:tinyint(4) unsigned;default:0;comment:消费类型 1：购买课时包 2：预约课程 3：无责取消预约 4：有责取消预约;NOT NULL" json:"change_type"`
}

func (StudentPackageDetail) TableName() string {
	return StudentPackageDetailTableName
}

var StudentPackageDetailFields = struct {
	ID               string
	CreatedAt        string
	UpdatedAt        string
	StudentId        string
	StudentPackageId string
	CourseId         string
	Title            string
	Hours            string
	LeftHours        string
	Change           string
	ChangeType       string
}{
	ID:               StudentPackageDetailTableName + ".id",
	CreatedAt:        StudentPackageDetailTableName + ".created_at",
	UpdatedAt:        StudentPackageDetailTableName + ".updated_at",
	StudentId:        StudentPackageDetailTableName + ".student_id",
	StudentPackageId: StudentPackageDetailTableName + ".student_package_id",
	CourseId:         StudentPackageDetailTableName + ".course_id",
	Title:            StudentPackageDetailTableName + ".title",
	Hours:            StudentPackageDetailTableName + ".hours",
	LeftHours:        StudentPackageDetailTableName + ".left_hours",
	Change:           StudentPackageDetailTableName + ".change",
	ChangeType:       StudentPackageDetailTableName + ".change_type",
}
