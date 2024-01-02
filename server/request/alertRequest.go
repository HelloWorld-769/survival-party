package request

import (
	"encoding/json"
)

type EmailRequest struct {
	ToEmail     string       `json:"toEmail" validate:"required"`
	Content     string       `json:"content" validate:"required"`
	Subject     string       `json:"subject" validate:"omitempty"`
	ContentType string       `json:"contentType" validate:"required"`
	Attachments []Attachment `json:"attachments,omitempty" validate:"omitempty"`
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
