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
	"golang.org/x/crypto/bcrypt"
)

// Helper function to create a user and get JWT token
func createUserAndGetToken(t *testing.T, userID, password, name string) string {
	ClearAll()
	
	// Register user
	requestBody := model.RegisterUserRequest{
		ID:       userID,
		Password: password,
		Name:     name,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Login to get token
	loginBody := model.LoginUserRequest{
		ID:       userID,
		Password: password,
	}

	loginJson, err := json.Marshal(loginBody)
	assert.Nil(t, err)

	loginReq := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(loginJson)))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Accept", "application/json")

	loginRes, err := app.Test(loginReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, loginRes.StatusCode)

	loginBytes, err := io.ReadAll(loginRes.Body)
	assert.Nil(t, err)

	loginResponse := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(loginBytes, loginResponse)
	assert.Nil(t, err)

	return loginResponse.Data.Token
}

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID:       "zhaka",
		Password: "rahasia",
		Name:     "Zhaka Hidayat",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID:       "",
		Password: "",
		Name:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	
	// Register first user
	firstRequestBody := model.RegisterUserRequest{
		ID:       "zhaka",
		Password: "rahasia",
		Name:     "Zhaka Hidayat",
	}

	firstBodyJson, err := json.Marshal(firstRequestBody)
	assert.Nil(t, err)

	firstRequest := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(firstBodyJson)))
	firstRequest.Header.Set("Content-Type", "application/json")
	firstRequest.Header.Set("Accept", "application/json")

	firstResponse, err := app.Test(firstRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, firstResponse.StatusCode)

	// Try to register duplicate user
	duplicateRequestBody := model.RegisterUserRequest{
		ID:       "zhaka", // Same ID
		Password: "rahasia123",
		Name:     "Another Name",
	}

	duplicateBodyJson, err := json.Marshal(duplicateRequestBody)
	assert.Nil(t, err)

	duplicateRequest := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(duplicateBodyJson)))
	duplicateRequest.Header.Set("Content-Type", "application/json")
	duplicateRequest.Header.Set("Accept", "application/json")

	duplicateResponse, err := app.Test(duplicateRequest)
	assert.Nil(t, err)

	duplicateBytes, err := io.ReadAll(duplicateResponse.Body)
	assert.Nil(t, err)

	duplicateResponseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(duplicateBytes, duplicateResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, duplicateResponse.StatusCode)
	assert.NotNil(t, duplicateResponseBody.Errors)
}

func TestLogin(t *testing.T) {
	ClearAll()
	
	// Register user first
	registerBody := model.RegisterUserRequest{
		ID:       "zhaka",
		Password: "rahasia",
		Name:     "Zhaka Hidayat",
	}

	registerJson, err := json.Marshal(registerBody)
	assert.Nil(t, err)

	registerReq := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(registerJson)))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Accept", "application/json")

	registerRes, err := app.Test(registerReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, registerRes.StatusCode)

	// Login
	requestBody := model.LoginUserRequest{
		ID:       "zhaka",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)
	assert.NotEmpty(t, responseBody.Data.Token)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
}

func TestLoginWrongUsername(t *testing.T) {
	ClearAll()
	
	// Register user first
	registerBody := model.RegisterUserRequest{
		ID:       "zhaka",
		Password: "rahasia",
		Name:     "Zhaka Hidayat",
	}

	registerJson, err := json.Marshal(registerBody)
	assert.Nil(t, err)

	registerReq := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(registerJson)))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Accept", "application/json")

	registerRes, err := app.Test(registerReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, registerRes.StatusCode)

	// Login with wrong username
	requestBody := model.LoginUserRequest{
		ID:       "wrong",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	
	// Register user first
	registerBody := model.RegisterUserRequest{
		ID:       "zhaka",
		Password: "rahasia",
		Name:     "Zhaka Hidayat",
	}

	registerJson, err := json.Marshal(registerBody)
	assert.Nil(t, err)

	registerReq := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(registerJson)))
	registerReq.Header.Set("Content-Type", "application/json")
	registerReq.Header.Set("Accept", "application/json")

	registerRes, err := app.Test(registerReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, registerRes.StatusCode)

	// Login with wrong password
	requestBody := model.LoginUserRequest{
		ID:       "zhaka",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestLogout(t *testing.T) {
	// Create user and get JWT token
	token := createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestLogoutWrongAuthorization(t *testing.T) {
	// Create user and get JWT token (but don't use it)
	createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestGetCurrentUser(t *testing.T) {
	// Create user and get JWT token
	token := createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "zhaka", responseBody.Data.ID)
	assert.Equal(t, "Zhaka Hidayat", responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	// Create user and get JWT token (but don't use it)
	createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestGetCurrentUserNoAuthorization(t *testing.T) {
	// Create user and get JWT token (but don't use it)
	createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// No Authorization header

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateUserName(t *testing.T) {
	// Create user and get JWT token
	token := createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	requestBody := model.UpdateUserRequest{
		Name: "Zhaka Hidayat Updated",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "zhaka", responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	// Create user and get JWT token
	token := createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "zhaka", responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// Verify password was actually updated in database
	user := new(entity.User)
	err = db.Where("id = ?", "zhaka").First(user).Error
	assert.Nil(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	assert.Nil(t, err)
}

func TestUpdateUserNameAndPassword(t *testing.T) {
	// Create user and get JWT token
	token := createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	requestBody := model.UpdateUserRequest{
		Name:     "Zhaka Hidayat Updated",
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "zhaka", responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// Verify password was actually updated in database
	user := new(entity.User)
	err = db.Where("id = ?", "zhaka").First(user).Error
	assert.Nil(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	assert.Nil(t, err)
}

func TestUpdateFailed(t *testing.T) {
	// Create user and get JWT token (but don't use it)
	createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	requestBody := model.UpdateUserRequest{
		Name:     "Should Fail",
		Password: "shouldfail",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateFailedNoAuthorization(t *testing.T) {
	// Create user and get JWT token (but don't use it)
	createUserAndGetToken(t, "zhaka", "rahasia", "Zhaka Hidayat")

	requestBody := model.UpdateUserRequest{
		Name:     "Should Fail",
		Password: "shouldfail",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// No Authorization header

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}