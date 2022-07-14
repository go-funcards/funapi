package session

import (
	"github.com/go-funcards/funapi/proto/user_service/v1"
)

type CreateUserDTO struct {
	Name           string `json:"name" validate:"required,min=3,max=100"`
	Email          string `json:"email" validate:"required,email,min=3,max=180"`
	Password       string `json:"password" validate:"required,eqfield=RepeatPassword,min=8,max=64"`
	RepeatPassword string `json:"repeat_password" validate:"required,min=8,max=64"`
}

func (dto CreateUserDTO) toCreate(id string, roles ...string) *v1.CreateUserRequest {
	return &v1.CreateUserRequest{
		UserId:   id,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Roles:    roles,
	}
}

type CredentialsDTO struct {
	Email    string `json:"email" validate:"required,email,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (dto CredentialsDTO) toRequest() *v1.UserByEmailAndPasswordRequest {
	return &v1.UserByEmailAndPasswordRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required,uuid4"`
}
