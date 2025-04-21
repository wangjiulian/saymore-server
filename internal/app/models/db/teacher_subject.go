package db

const TeacherSubjectTableName = "teacher_subjects"

// TeacherSubject 老师授课科目关联表
type TeacherSubject struct {
	Base
	TeacherId int64 `gorm:"column:teacher_id;type:bigint(20);default:0;comment:老师ID;NOT NULL" json:"teacher_id"`
	SubjectId int   `gorm:"column:subject_id;type:int(11);default:0;comment:科目ID;NOT NULL" json:"subject_id"`
}

func (TeacherSubject) TableName() string {
	return TeacherSubjectTableName
}

var TeacherSubjectFields = struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	TeacherId string
	SubjectId string
}{
	ID:        TeacherSubjectTableName + ".id",
	CreatedAt: TeacherSubjectTableName + ".created_at",
	UpdatedAt: TeacherSubjectTableName + ".updated_at",
	TeacherId: TeacherSubjectTableName + ".teacher_id",
	SubjectId: TeacherSubjectTableName + ".subject_id",
}
