package server

import (
	"encoding/json"
	"fmt"
	"log"
	_ "main/docs"
	"main/server/handler"
	gomail "main/server/services/alert_service/Gomail"
	"main/server/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateUser(ctx *gin.Context) {

	var user request.CreateUser
	utils.RequestDecoding(ctx, &user)

	exUser := model.UserInfo{Hobby: "cricket", Category: "admin"}

	exUserByte, err := json.Marshal(exUser)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("exuserDATA", string(exUserByte))
	return

	// userNew := model.User{

	// 	Email:    "test@example.com",
	// 	Password: "password123",
	// 	FullName: "John Doe",
	// 	Info: json.RawMessage(`
	// 		{
	// 			"hobby": "Reading",
	// 			"category": "Bookworm"
	// 		}
	// 	`),
	// }

	// err = db.CreateRecord(&userNew)
	// if err != nil {

	// 	log.Fatal(err)
	// 	return
	// }

	// fmt.Println("user created successfully")
	// response.ShowResponse("user created successfully", utils.HTTP_OK, "success", nil, ctx)
}

func ConfigureRoutes(server *Server) {

	//Allowing CORS
	server.engine.Use(gateway.CORSMiddleware())

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/signup", handler.SignupHandler)
	server.engine.POST("/login", handler.LoginHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/send-email", gomail.SendEmailOtpService)
	server.engine.POST("/send-sms", handler.TwilioServiceHnadler)
	server.engine.POST("/create-user", CreateUser)

	server.engine.GET("/ping", Pong)

}

// For server Testing(acknowledgement)
func Pong(ctx *gin.Context) {

	msg := "Pong"
	ctx.Writer.Write([]byte(msg))
}
