package db

const TeacherAvailabilityTableName = "teacher_availabilities"

// TeacherAvailability 老师可预约时间表
type TeacherAvailability struct {
	Base
	TeacherId int64  `gorm:"column:teacher_id;type:bigint(20);default:0;comment:关联teacher表的ID;NOT NULL" json:"teacher_id"`
	CourseId  int64  `gorm:"column:course_id;type:bigint(20);default:0;comment:关联courese表的ID;NOT NULL" json:"course_id"`
	StartTime uint64 `gorm:"column:start_time;type:bigint(20);comment:可预约开始时间（Unix时间戳，单位：秒）;NOT NULL" json:"start_time"`
	EndTime   uint64 `gorm:"column:end_time;type:bigint(20);comment:可预约结束时间（Unix时间戳，单位：秒）;NOT NULL" json:"end_time"`
}

func (TeacherAvailability) TableName() string {
	return TeacherAvailabilityTableName
}

var TeacherAvailabilityFields = struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	TeacherId string
	CourseId  string
	StartTime string
	EndTime   string
}{
	ID:        TeacherAvailabilityTableName + ".id",
	CreatedAt: TeacherAvailabilityTableName + ".created_at",
	UpdatedAt: TeacherAvailabilityTableName + ".updated_at",
	TeacherId: TeacherAvailabilityTableName + ".teacher_id",
	CourseId:  TeacherAvailabilityTableName + ".course_id",
	StartTime: TeacherAvailabilityTableName + ".start_time",
	EndTime:   TeacherAvailabilityTableName + ".end_time",
}
