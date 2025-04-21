package db

const FeedbackTableName = "feedbacks"

// Feedbacks 意见反馈表
type Feedbacks struct {
	Base
	StudentId int64  `gorm:"column:student_id;type:bigint(20);default:0;comment:关联学生表的 ID;NOT NULL" json:"student_id"`
	Content   string `gorm:"column:content;type:varchar(255);comment:反馈内容;NOT NULL" json:"content"`
}

func (Feedbacks) TableName() string {
	return FeedbackTableName
}

var FeedbackFields = struct {
	ID      string
	Student string
	Content string
}{
	ID:      FeedbackTableName + ".id",
	Student: FeedbackTableName + ".student_id",
	Content: FeedbackTableName + ".content",
}
