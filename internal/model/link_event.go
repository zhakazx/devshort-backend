package model

type LinkEvent struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Title     string `json:"title"`
	ShortUrl  string `json:"short_url"`
	LongUrl   string `json:"long_url"`
	IsActive  bool   `json:"is_active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (l *LinkEvent) GetId() string {
	return l.ID
}