package gateway

import (
	// "fmt"
	// "gym/server/response"

	"fmt"
	"main/server/db"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthorization(ctx *gin.Context) {

	fmt.Println("inside middleware")
	bearerToken := ctx.Request.Header.Get("Authorization")

	tokenString := strings.Split(bearerToken, " ")[1]

	var exists bool
	//first check if the session is valid or not
	query := "SELECT EXISTS(SELECT 1 FROM sessions WHERE token=?)"
	err := db.QueryExecutor(query, &exists, tokenString)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}
	if !exists {
		response.ShowResponse("Invalid session", utils.HTTP_FORBIDDEN, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}

	claims, err := token.DecodeToken(tokenString)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}
	fmt.Println("claims:", claims)
	err = claims.Valid()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}

	ctx.Set("userId", claims.Id)

	ctx.Next()

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
