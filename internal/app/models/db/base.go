package db

type Base struct {
	Id        int64 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt int64 `gorm:"column:created_at;type:int;default:0;comment:Creation time;NOT NULL" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:int;default:0;comment:Update time;NOT NULL" json:"updated_at"`
}
