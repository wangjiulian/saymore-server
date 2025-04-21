package db

const SubjectTableName = "subjects"

// Subject 科目表
type Subject struct {
	Base
	Name      string `gorm:"column:name;type:varchar(32);comment:名称;NOT NULL" json:"name"`
	Code      string `gorm:"column:code;type:varchar(32);comment:编码;NOT NULL" json:"code"`
	ParentId  int    `gorm:"column:parent_id;type:int(11);default:0;comment:父级ID，顶级为NULL" json:"parent_id"`
	SortOrder int    `gorm:"column:sort_order;type:int(11);default:0;comment:排序字段;NOT NULL" json:"sort_order"`
}

func (Subject) TableName() string {
	return SubjectTableName
}

var SubjectFields = struct {
	ID        string
	Name      string
	Code      string
	ParentId  string
	SortOrder string
	CreatedAt string
	UpdatedAt string
}{
	ID:        SubjectTableName + ".id",
	Name:      SubjectTableName + ".name",
	Code:      SubjectTableName + ".code",
	ParentId:  SubjectTableName + ".parent_id",
	SortOrder: SubjectTableName + ".sort_order",
	CreatedAt: SubjectTableName + ".created_at",
	UpdatedAt: SubjectTableName + ".updated_at",
}
