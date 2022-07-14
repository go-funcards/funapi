package users

import (
	"github.com/go-funcards/funapi/proto/user_service/v1"
	"time"
)

type UpdateUserDTO struct {
	UserID            string `json:"-" uri:"user_id" validate:"required,uuid4"`
	Name              string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Email             string `json:"email,omitempty" validate:"omitempty,email,min=3,max=255"`
	OldPassword       string `json:"old_password,omitempty" validate:"omitempty,required_with=NewPassword,min=8,max=64"`
	NewPassword       string `json:"new_password,omitempty" validate:"omitempty,eqfield=RepeatNewPassword,min=8,max=64"`
	RepeatNewPassword string `json:"repeat_new_password,omitempty" validate:"omitempty,min=8,max=64"`
}

func (dto UpdateUserDTO) toUpdate() *v1.UpdateUserRequest {
	return &v1.UpdateUserRequest{
		UserId:      dto.UserID,
		Name:        dto.Name,
		Email:       dto.Email,
		OldPassword: dto.OldPassword,
		NewPassword: dto.NewPassword,
	}
}

type User struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(response *v1.UserResponse) User {
	return User{
		UserID:    response.GetUserId(),
		Name:      response.GetName(),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
	}
}
