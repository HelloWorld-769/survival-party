package response

import (
	"github.com/gin-gonic/gin"
)

type Success struct {
	Status Status `json:"status"`
}
type Status struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(context *gin.Context, statusCode int, data interface{}) {
	context.JSON(statusCode, data)
}

func ShowResponse(message string, statusCode int64, status string, data interface{}, context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Writer.WriteHeader(int(statusCode))
	response := Success{
		Status: Status{
			Code:    statusCode,
			Message: message,
			Data:    data,
		},
	}

	Response(context, int(statusCode), response)
}
