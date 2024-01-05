package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func RequestDecoding(ctx *gin.Context, data interface{}) error {

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	fmt.Println("reqBody        ", string(reqBody))
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		return err
	}
	return nil
}

func SetHeader(context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "application/json")

}
