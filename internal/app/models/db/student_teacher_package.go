package db

const StudentTeacherPackageTableName = "student_teacher_packages"

// StudentTeacherPackages 学生老师课时包关联表
type StudentTeacherPackages struct {
	Base
	TeacherId        int64   `gorm:"column:teacher_id;type:bigint(20);default:0;comment:教师ID;NOT NULL" json:"teacher_id"`
	StudentId        int64   `gorm:"column:student_id;type:bigint(20);default:0;comment:学生ID;NOT NULL" json:"student_id"`
	StudentPackageId int64   `gorm:"column:student_package_id;type:bigint(20);default:0;comment:student_packages表ID;NOT NULL" json:"student_package_id"`
	Price            float64 `gorm:"column:price;type:decimal(8,2);default:0.00;comment:课时费;NOT NULL" json:"price"`
}

func (StudentTeacherPackages) TableName() string {
	return StudentTeacherPackageTableName
}

var StudentTeacherPackageFields = struct {
	ID               string
	CreatedAt        string
	UpdatedAt        string
	TeacherId        string
	StudentId        string
	StudentPackageId string
	Price            string
}{
	ID:               StudentTeacherPackageTableName + ".id",
	CreatedAt:        StudentTeacherPackageTableName + ".created_at",
	UpdatedAt:        StudentTeacherPackageTableName + ".updated_at",
	TeacherId:        StudentTeacherPackageTableName + ".teacher_id",
	StudentId:        StudentTeacherPackageTableName + ".student_id",
	StudentPackageId: StudentTeacherPackageTableName + ".student_package_id",
	Price:            StudentTeacherPackageTableName + ".price",
}
