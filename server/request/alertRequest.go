package request 



type EmailRequest struct {
	ToEmail     string       `json:"toEmail" validate:"required"`
	Content     string       `json:"content" validate:"required"`
	Subject     string       `json:"subject" validate:"omitempty"`
	ContentType string       `json:"contentType" validate:"required"`
	Attachments []Attachment `json:"attachments,omitempty" validate:"omitempty"`
}

type Attachment struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}


type TwilioSmsRequest struct {

	Contact string `json:"contact" validate:"required`
}
