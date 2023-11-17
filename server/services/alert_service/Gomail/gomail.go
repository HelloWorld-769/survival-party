package gomail

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"main/server/validation"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
)

// func SendEmailOtpService(context *gin.Context, req request.EmailRequest) {

// 	//get the content type from the header

// 	contentType:=context.Request.Header.Get("Content-Type")
// 	m := gomail.NewMessage()

// 	// Set E-Mail sender
// 	m.SetHeader("From", "amantarar01@gmail.com")

// 	// Set E-Mail receivers
// 	m.SetHeader("To", req.ToEmail)

// 	// Set E-Mail subject
// 	// m.SetHeader("Subject", "Instagram Email verification")

// 	// Set E-Mail body. You can set plain text or html with text/html
// 	// rand.Seed(time.Now().UnixNano())
// 	// otp := rand.Int()
// 	// strOtp := strconv.Itoa(otp)
// 	m.SetBody(contentType, req.Content)

// 	// Settings for SMTP server
// 	d := gomail.NewDialer("smtp.gmail.com", 587, "amantarar01@gmail.com", "mdyrprmdvxpfxjjp")

// 	// This is only needed when SSL/TLS certificate is not valid on server.
// 	// In production this should be set to false.
// 	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 	// Now send E-Mail
// 	if err := d.DialAndSend(m); err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	// var usersession model.UserAuthSessions

// 	// usersession.Email = toEmail
// 	// er := db.CreateRecord(&usersession)
// 	// if er != nil {
// 	// 	response.ShowResponse("Db ERROR", 500, er.Error(), "", context)
// 	// 	return
// 	// }
// 	// //set a sessionID cookie
// 	// db.FindById(&usersession, toEmail, "email")

// 	// // sessionIdCookie := &http.Cookie{Name: "sessionID", Value: usersession.SessionId, Domain: ".rr.com"}
// 	// // http.SetCookie(context.Writer, sessionIdCookie)
// 	// // fmt.Println("session cookie is set")

// 	// //set a cookie with hash value of otp
// 	// hash, _ := utils.HashPassword(strOtp)
// 	// // cookie := &http.Cookie{Name: "otpHash", Value: hash, SameSite: 4}
// 	// // http.SetCookie(context.Writer, cookie)

// 	claims := model.Claims{Email: toEmail, SessionId: usersession.SessionId, Hash: hash}
// 	token := utils.GenerateToken(claims)
// 	context.Writer.Header().Set("UserToken", token)

// 	response.ShowResponse("Success", 200, "Code sent on Email", token, context)

// }

func isHTML(content string) bool {
	// Check if the content contains HTML tags
	return strings.Contains(content, "<html>") || strings.Contains(content, "<body>")
}

// isXML checks if the provided string contains XML content.
func isXML(content string) bool {
	var parsedData interface{}
	err := xml.Unmarshal([]byte(content), &parsedData)
	return err == nil
}

func SendEmailOtpService(context *gin.Context) {

	utils.SetHeader(context)

	var req request.EmailRequest

	utils.RequestDecoding(context, &req)

	//validation Check on request body fields
	err := validation.CheckValidation(&req)
	if err != nil {
		response.ShowResponse(err.Error(), 400, "Failure", "", context)
		return
	}

	//check the content type
	if !isHTML(req.Content) && !isXML(req.Content) {
		if req.ContentType != "text/plain" {
			//show error
			response.ShowResponse("Enter correct content type as text/plain", 400, "Failure", "", context)
			return
		}
	}

	if isHTML(req.Content) && req.ContentType != "text/html" {

		//show error
		response.ShowResponse("Enter correct content type as text/html", 400, "Failure", "", context)
		return
	}
	if isXML(req.Content) && req.ContentType != "application/xml" {

		//show error
		response.ShowResponse("Enter correct content type as application/xml content", 400, "Failure", "", context)
		return

	}

	// Create a new message
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("FROM_EMAIL"))

	// contentType := context.Request.Header.Get("Content-Type")

	// Set E-Mail receivers
	m.SetHeader("To", req.ToEmail)

	// Set E-Mail subject
	m.SetHeader("Subject", req.Subject)

	// Set E-Mail body. You can set plain text or html with text/html

	m.SetBody(req.ContentType, req.Content)

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
