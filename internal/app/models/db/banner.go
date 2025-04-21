package db

const BannerTableName = "banners"

type Banner struct {
	Base
	Title     string `gorm:"column:title;type:varchar(32);comment:Title;NOT NULL" json:"title"`
	Image     string `gorm:"column:image;type:varchar(255);comment:Image URL;NOT NULL" json:"image"`
	Url       string `gorm:"column:url;type:varchar(255);comment:Link URL;NOT NULL" json:"url"`
	SortOrder int    `gorm:"column:sort_order;type:int(11);default:0;comment:Sort order;NOT NULL" json:"sort_order"`
}

func (Banner) TableName() string {
	return BannerTableName
}

var BannerFields = struct {
	ID        string
	Title     string
	Image     string
	Url       string
	SortOrder string
}{
	ID:        BannerTableName + ".id",
	Title:     BannerTableName + ".title",
	Image:     BannerTableName + ".image",
	Url:       BannerTableName + ".url",
	SortOrder: BannerTableName + ".sort_order",
}
