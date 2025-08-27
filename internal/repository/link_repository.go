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

func (r *LinkRepository) FindById(tx *gorm.DB, id string) (*entity.Link, error) {
	link := new(entity.Link)
	if err := tx.Where("id = ?", id).First(link).Error; err != nil {
		r.Log.WithError(err).Error("error finding link by id")
		return nil, err
	}
	return link, nil
}

func (r *LinkRepository) FindAllByUserId(tx *gorm.DB, userId string) ([]entity.Link, error) {
	var links []entity.Link
	if err := tx.Where("user_id = ?", userId).Find(&links).Error; err != nil {
		r.Log.WithError(err).Error("error finding links by user id")
		return nil, err
	}
	return links, nil
}

func (r *LinkRepository) FindActiveLinksByUserId(tx *gorm.DB, userId string) ([]entity.Link, error) {
    var links []entity.Link
    if err := tx.Where("user_id = ? AND is_active = ?", userId, true).Find(&links).Error; err != nil {
        r.Log.WithError(err).Error("error finding active links by user id")
        return nil, err
    }
    return links, nil
}
