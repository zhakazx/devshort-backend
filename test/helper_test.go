package test

import (
	"devshort-backend/internal/entity"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ClearAll() {
	ClearUsers()
	ClearLinks()
}

func ClearUsers() {
	err := db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func ClearLinks() {
	err := db.Where("id is not null").Delete(&entity.Link{}).Error
	if err != nil {
		log.Fatalf("Failed clear link data : %+v", err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}

func GetFirstLink(t *testing.T, user *entity.User) *entity.Link {
	link := new(entity.Link)
	err := db.Where("user_id = ?", user.ID).First(link).Error
	assert.Nil(t, err)
	return link
}

func CreateLinks(user *entity.User, total int) {
	for i := 0; i < total; i++ {
		link := &entity.Link{
			ID:       uuid.NewString(),
			UserId:   user.ID,
			Title:    "Link " + strconv.Itoa(i),
			ShortUrl: "link" + strconv.Itoa(i),
			LongUrl:  "https://example" + strconv.Itoa(i) + ".com",
			IsActive: i%2 == 0, // Alternate active status
		}
		err := db.Create(link).Error
		if err != nil {
			log.Fatalf("Failed create link data : %+v", err)
		}
	}
}
