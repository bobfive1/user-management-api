package userprofile

import (
	"time"

	"github.com/bobfive1/user-management-api/internal/chrono"
)

type UserProfile struct {
	UserID    string           `db:"user_id" json:"user_id"`
	Password  string           `db:"password" json:"-"`
	FirstName string           `db:"first_name" json:"first_name"`
	LastName  string           `db:"last_name" json:"last_name"`
	Address   *string          `db:"address" json:"address"`
	Birthdate *chrono.DateOnly `db:"birthdate" json:"birthdate"`
	Email     *string          `db:"email" json:"email"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" json:"updated_at"`
}

type UserProfileDisplay struct {
	UserID    string           `db:"user_id" json:"user_id"`
	FirstName string           `db:"first_name" json:"first_name"`
	LastName  string           `db:"last_name" json:"last_name"`
	Address   *string          `db:"address" json:"address"`
	Birthdate *chrono.DateOnly `db:"birthdate" json:"birthdate"`
	Email     *string          `db:"email" json:"email"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt time.Time        `db:"updated_at" json:"updated_at"`
}

func NewUserProfileDisplay(user UserProfile) UserProfileDisplay {
	return UserProfileDisplay{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Address:   user.Address,
		Birthdate: user.Birthdate,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type InsertUserProfileRequest struct {
	UserID    string           `db:"user_id" json:"user_id" binding:"required,max=20"`
	Password  string           `db:"password" json:"password" binding:"required,max=20"`
	FirstName string           `json:"first_name" binding:"required,max=150"`
	LastName  string           `json:"last_name" binding:"required,max=150"`
	Address   *string          `json:"address"`
	Birthdate *chrono.DateOnly `json:"birthdate" binding:"omitempty,checkyear"`
	Email     *string          `json:"email" binding:"omitempty,email,max=255"`
}

type UpdateUserProfileRequest struct {
	Password  string           `db:"password" json:"password" binding:"required,max=20"`
	FirstName string           `json:"first_name" binding:"required,max=150"`
	LastName  string           `json:"last_name" binding:"required,max=150"`
	Address   *string          `json:"address"`
	Birthdate *chrono.DateOnly `json:"birthdate" binding:"omitempty,checkyear"`
	Email     *string          `json:"email" binding:"omitempty,email,max=255"`
}

type UserProfileResponse struct {
	TimeStamp    time.Time `json:"timestamp"`
	ErrorCode    string    `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	ErrorDetail  any       `json:"error_detail,omitempty"`
	Data         any       `json:"data,omitempty"`
}

type UserProfileLoginRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUserProfileResponse(code, message string, data any) UserProfileResponse {
	return UserProfileResponse{
		TimeStamp:    time.Now(),
		ErrorCode:    code,
		ErrorMessage: message,
		Data:         data,
	}

}
