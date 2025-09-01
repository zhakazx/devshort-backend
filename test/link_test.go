package test

import (
	"devshort-backend/internal/entity"
	"devshort-backend/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLink(t *testing.T) {
	TestLogin(t)

	user := new(entity.User)
	err := db.Where("id = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.CreateLinkRequest{
		UserId:   user.ID,
		Title:    "Test Link",
		ShortUrl: "test123",
		LongUrl:  "https://example.com",
		IsActive: true,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[*model.LinkResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Title, responseBody.Data.Title)
	assert.Equal(t, requestBody.ShortUrl, responseBody.Data.ShortUrl)
	assert.Equal(t, requestBody.LongUrl, responseBody.Data.LongUrl)
	assert.Equal(t, requestBody.IsActive, responseBody.Data.IsActive)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}