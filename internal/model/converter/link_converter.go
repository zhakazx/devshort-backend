package converter

import (
	"devshort-backend/internal/entity"
	"devshort-backend/internal/model"
)

func LinkToResponse(link *entity.Link) *model.LinkResponse {
	return &model.LinkResponse{
		ID:       link.ID,
		Title:    link.Title,
		ShortUrl: link.ShortUrl,
		LongUrl:  link.LongUrl,
		IsActive: link.IsActive,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
	}
}

func LinkToEvent(link *entity.Link) *model.LinkEvent {
	return &model.LinkEvent{
		ID:       link.ID,
		UserId:   link.UserId,
		Title:    link.Title,
		ShortUrl: link.ShortUrl,
		LongUrl:  link.LongUrl,
		IsActive: link.IsActive,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
	}
}