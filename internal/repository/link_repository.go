package repository

import (
	"devshort-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LinkRepository struct {
	Repository[entity.Link]
	Log *logrus.Logger
}

func NewLinkRepository(log *logrus.Logger) *LinkRepository {
	return &LinkRepository{
		Log: log,
	}
}

func (r *LinkRepository) FindByIdAndUserId(tx *gorm.DB, link *entity.Link, id string, userId string) error {
	return tx.Where("id = ? AND user_id = ?", id, userId).First(link).Error
}

func (r *LinkRepository) FindAllByUserId(tx *gorm.DB, userId string) ([]entity.Link, error) {
	var links []entity.Link
	if err := tx.Where("user_id = ?", userId).Find(&links).Error; err != nil {
		r.Log.WithError(err).Error("error finding links by user id")
		return nil, err
	}
	return links, nil
}
