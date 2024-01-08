package Gomail

import (
	"crypto/tls"
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
)

// generateOTP generates a 6-digit OTP
func generateOTP() int {
	rand.Seed(time.Now().UnixNano())

	// Generate a random 6-digit number
	otp := rand.Intn(1000000)

	// Ensure the OTP has exactly 6 digits
	otp = otp % 1000000

	return otp
}

func SendEmailOtpService(context *gin.Context, req request.EmailRequest) {

	utils.SetHeader(context)

	// Create a new message
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("FROM_EMAIL"))

	// contentType := context.Request.Header.Get("Content-Type")

	// Set E-Mail receivers
	m.SetHeader("To", req.Users.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "SURVIVAL RESET PASSWORD OTP")

	// Set E-Mail body. You can set plain text or html with text/html

	//generate random 6 digit OTP
	otp := generateOTP()
	m.SetBody("text/plain", fmt.Sprint(otp))
	// m.Attach("/home/chicmic/Downloads/image.png")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("FROM_EMAIL"), os.Getenv("EMAIL_PASS"))

	// This is only needed when SSL/TLS certificate is not valid on the server.
	// In production, this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	response.ShowResponse("email sent successfully", 200, "Success", "", context)
}
