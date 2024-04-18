package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SigupRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Avatar   int64  `json:"avatar"`
	} `json:"user"`
}

func (a SigupRequest) Validate() error {
	return validation.ValidateStruct(&a.User,
		validation.Field(&a.User.Username, validation.Required),
		validation.Field(&a.User.Password, validation.Required),
		validation.Field(&a.User.Email, validation.Required, is.Email),
		validation.Field(&a.User.Avatar, validation.Required),
	)
}

type LoginRequest struct {
	User struct {
		Credential string `json:"credential"`
		Password   string `json:"password"`
	} `json:"user"`
}

func (a LoginRequest) Validate() error {
	return validation.ValidateStruct(&a.User,
		validation.Field(&a.User.Credential, validation.Required),
		validation.Field(&a.User.Password, validation.Required),
	)
}

type SocialLoginReq struct {
	Email  string `json:"email"`
	Avatar int64  `json:"avatar"`
	Uid    string `json:"uid"`
}

func (a SocialLoginReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Avatar, validation.Required),
		validation.Field(&a.Uid, validation.Required),
		validation.Field(&a.Email, validation.Required, is.Email),
	)
}
