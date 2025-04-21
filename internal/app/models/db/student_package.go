package db

import "github.com/shopspring/decimal"

const StudentPackageTableName = "student_packages"

// StudentPackage 学生课时包表
type StudentPackage struct {
	Base
	StudentId   int64           `gorm:"column:student_id;type:bigint(20);default:0;comment:关联student表的ID;NOT NULL" json:"student_id"`
	Name        string          `gorm:"column:name;type:varchar(100);comment:课时包名称;NOT NULL" json:"name"`
	SubjectId   int             `gorm:"column:subject_id;type:int(11);default:0;comment:关联subjects表的一级父类ID;NOT NULL" json:"subject_id"`
	Hours       decimal.Decimal `gorm:"column:hours;type:decimal(6,2) unsigned;default:0.0;comment:课时数量;NOT NULL" json:"hours"`
	LeftHours   decimal.Decimal `gorm:"column:left_hours;type:decimal(6,2) unsigned;default:0.0;comment:剩余课时包数量;NOT NULL" json:"left_hours"`
	Description string          `gorm:"column:description;type:varchar(255);comment:课时包说明;NOT NULL" json:"description"`
}

func (StudentPackage) TableName() string {
	return StudentPackageTableName
}

var StudentPackageFields = struct {
	ID          string
	CreatedAt   string
	UpdatedAt   string
	StudentId   string
	Name        string
	SubjectId   string
	Hours       string
	LeftHours   string
	Description string
}{
	ID:          StudentPackageTableName + ".id",
	CreatedAt:   StudentPackageTableName + ".created_at",
	UpdatedAt:   StudentPackageTableName + ".updated_at",
	StudentId:   StudentPackageTableName + ".student_id",
	Name:        StudentPackageTableName + ".name",
	SubjectId:   StudentPackageTableName + ".subject_id",
	Hours:       StudentPackageTableName + ".hours",
	LeftHours:   StudentPackageTableName + ".left_hours",
	Description: StudentPackageTableName + ".description",
}
