package db

const StudentTrialQuotaTableName = "student_trial_quotas"

// StudentTrialQuota 学生试听额度表
type StudentTrialQuota struct {
	Base
	StudentId  int64 `gorm:"column:student_id;type:bigint(20);default:0;comment:关联student表的ID;NOT NULL" json:"student_id"`
	CourseId   int64 `gorm:"column:course_id;type:bigint(20);default:0;comment:关联courese表的ID;NOT NULL" json:"course_id"`
	OperatorId int64 `gorm:"column:operator_id;type:bigint(20);default:0;comment:操作人ID(教务老师);NOT NULL" json:"operator_id"`
}

func (StudentTrialQuota) TableName() string {
	return StudentTrialQuotaTableName
}

var StudentTrialQuotaFields = struct {
	ID         string
	CreatedAt  string
	UpdatedAt  string
	StudentId  string
	CourseId   string
	OperatorId string
}{
	ID:         StudentTrialQuotaTableName + ".id",
	CreatedAt:  StudentTrialQuotaTableName + ".created_at",
	UpdatedAt:  StudentTrialQuotaTableName + ".updated_at",
	StudentId:  StudentTrialQuotaTableName + ".student_id",
	CourseId:   StudentTrialQuotaTableName + ".course_id",
	OperatorId: StudentTrialQuotaTableName + ".operator_id",
}
