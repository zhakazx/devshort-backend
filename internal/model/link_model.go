package model

type LinkResponse struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Title     string `json:"title"`
	ShortUrl  string `json:"short_url"`
	LongUrl   string `json:"long_url"`
	IsActive  bool   `json:"is_active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type ListLinkRequest struct {
	UserId string `json:"-" validate:"required,max=100,uuid"`
}

type CreateLinkRequest struct {
	UserId   string `json:"-" validate:"required,max=100,uuid"`
	Title    string `json:"title" validate:"required,min=2,max=50"`
	ShortUrl string `json:"short_url" validate:"required"`
	LongUrl  string `json:"long_url" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type GetLinkRequest struct {
	ID     string `json:"-" validate:"required,uuid"`
	UserId string `json:"-" validate:"required,max=100,uuid"`
}

type UpdateLinkRequest struct {
	ID       string `json:"-" validate:"required,uuid"`
	UserId   string `json:"-" validate:"required,max=100,uuid"`
	Title    string `json:"title" validate:"required,min=2,max=50"`
	ShortUrl string `json:"short_url" validate:"required"`
	LongUrl  string `json:"long_url" validate:"required"`
	IsActive bool   `json:"is_active" validate:"required"`
}

type DeleteLinkRequest struct {
	ID     string `json:"-" validate:"required,uuid"`
	UserId string `json:"-" validate:"required,max=100,uuid"`
}
