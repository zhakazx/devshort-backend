package usecase

import (
	"context"
	"devshort-backend/internal/entity"
	"devshort-backend/internal/gateway/messaging"
	"devshort-backend/internal/model"
	"devshort-backend/internal/model/converter"
	"devshort-backend/internal/repository"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	UserProducer   *messaging.UserProducer
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, userProducer *messaging.UserProducer) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		UserProducer:   userProducer,
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (c *UserUseCase) generateJWT(userID string) (string, error) {
    claims := jwt.MapClaims{
        "id":  userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(), // expired 1 hari
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
    token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fiber.ErrUnauthorized
        }
        return jwtSecret, nil
    })

    if err != nil || !token.Valid {
        return nil, fiber.ErrUnauthorized
    }

    claims := token.Claims.(jwt.MapClaims)
    return &model.Auth{ID: claims["id"].(string)}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	total, err := c.UserRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed count user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("User already exists : %+v", err)
		return nil, fiber.ErrConflict
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       request.ID,
		Password: string(password),
		Name:     request.Name,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user created event")
		if err = c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user created event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user created event")
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
    tx := c.DB.WithContext(ctx).Begin()
    defer tx.Rollback()

    if err := c.Validate.Struct(request); err != nil {
        c.Log.Warnf("Invalid request body  : %+v", err)
        return nil, fiber.ErrBadRequest
    }

    user := new(entity.User)
    if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
        c.Log.Warnf("Failed find user by id : %+v", err)
        return nil, fiber.ErrUnauthorized
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
        c.Log.Warnf("Invalid password : %+v", err)
        return nil, fiber.ErrUnauthorized
    }

    token, err := c.generateJWT(user.ID)
    if err != nil {
        c.Log.Warnf("Failed to generate jwt : %+v", err)
        return nil, fiber.ErrInternalServerError
    }

    if err := tx.Commit().Error; err != nil {
        c.Log.Warnf("Failed commit transaction : %+v", err)
        return nil, fiber.ErrInternalServerError
    }

    if c.UserProducer != nil {
        event := converter.UserToEvent(user)
        c.Log.Info("Publishing user login event")
        if err := c.UserProducer.Send(event); err != nil {
            c.Log.Warnf("Failed publish user login event : %+v", err)
            return nil, fiber.ErrInternalServerError
        }
    } else {
        c.Log.Info("Kafka producer is disabled, skipping user login event")
    }

    return &model.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Token: token,
    }, nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
    c.Log.Infof("User %s logged out (client-side token deletion)", request.ID)
    return true, nil
}


func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user updated event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user updated event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user updated event")
	}

	return converter.UserToResponse(user), nil
}
