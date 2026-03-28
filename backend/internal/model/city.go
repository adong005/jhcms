package model

type City struct {
	CityCode int    `gorm:"primaryKey;column:city_code" json:"cityCode"`
	Name     string `gorm:"type:varchar(64);index;not null" json:"name"`
	Pinyin   string `gorm:"type:varchar(64);index;not null" json:"pinyin"`
	Status   int8   `gorm:"default:1;index;not null" json:"status"`
	AuditModel
}

func (City) TableName() string {
	return "city_list"
}
