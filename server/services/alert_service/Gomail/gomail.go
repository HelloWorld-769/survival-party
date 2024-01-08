package Gomail

import (
	"crypto/tls"
	"fmt"
	"main/server/db"
	"main/server/model"
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

	//store this otp in RestSessions of user in databsae
	var userResetRession model.ResetSession

	//check if reset Session is already present for the user ,then update OTP only
	exists, err := ResetSessionAlreadyPresent(req.Users.Email)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, context)
		return
	}
	if exists {
		//update the otp in reset session only
		userResetRession.Otp = int64(otp)
		err := db.UpdateRecord(&userResetRession, req.Users.Email, "user_email").Error
		if err != nil {
			response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, context)
			return
		}

		response.ShowResponse("email sent successfully", 200, "Success", "", context)
		return

	}
	userResetRession.UserEmail = req.Users.Email
	userResetRession.Otp = int64(otp)
	err = db.CreateRecord(&userResetRession)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, context)
		return
	}

	response.ShowResponse("email sent successfully", 200, "Success", "", context)
}

func SendEmailService(context *gin.Context, link string, toEmail string) {

	utils.SetHeader(context)

	// Create a new message
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("FROM_EMAIL"))

	// contentType := context.Request.Header.Get("Content-Type")

	// Set E-Mail receivers
	m.SetHeader("To", toEmail)

	// Set E-Mail subject
	m.SetHeader("Subject", "SURVIVAL RESET PASSWORD OTP")

	// Set E-Mail body. You can set plain text or html with text/html

	m.SetBody("text/plain", link)
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

}
func ResetSessionAlreadyPresent(email string) (bool, error) {

	var exists bool = false
	query := "select exists(select * from reset_sessions where user_email=?)"
	err := db.QueryExecutor(query, &exists, email)
	if err != nil {
		return false, err
	}
	return exists, nil

}
