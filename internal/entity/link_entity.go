package entity

type Link struct {
	ID        string  `gorm:"column:id;primaryKey"`
	UserId    string  `gorm:"column:user_id"`
	Title     string  `gorm:"column:title"`
	ShortUrl  string  `gorm:"column:short_url"`
	LongUrl   string  `gorm:"column:long_url"`
	IsActive  bool    `gorm:"column:is_active"`
	CreatedAt int64   `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64   `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	User      User    `gorm:"foreignKey:user_id;references:id"`
}

func (a *Link) TableName() string {
	return "links"
}
