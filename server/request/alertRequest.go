package request

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type EmailRequest struct {
	Users struct {
		Email string `json:"email" validate:"required"`
	} `json:"user"`
	// Content     string `json:"content"`
	// Subject     string `json:"subject" omitempty`
	// ContentType string `json:"contentType"`
}

func (a *EmailRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(a.Users.Email, validation.Required, is.Email),
	)
}

type RestPasswordRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type OtpRequest struct {
	Email string `json:"email" validate:"required"`
	Otp   int64  `json:"otp" validate:"required"`
}

func (a *OtpRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(a.Email, validation.Required, is.Email),
		validation.Field(a.Otp, validation.Required),
	)
}

func (a *RestPasswordRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(a.Email, validation.Required, is.Email),
		validation.Field(a.Password, validation.Required),
	)
}

type CreateUser struct {
	Name string          `json:"name" validate:"required"`
	Age  int             `json:"age" validate:"required"`
	Info json.RawMessage `json:"user_info"`
}

type Attachment struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

type TwilioSmsRequest struct {
	Contact string `json:"contact" validate:"required"`
}
