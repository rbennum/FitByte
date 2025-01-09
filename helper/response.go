package helper

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func NewResponse(data interface{}, error error) *Response {

	if error != nil {
		return &Response{
			Data:  data,
			Error: error.Error(),
		}
	}

	return &Response{
		data,
		error,
	}
}

func FallbackResponse(ctx *gin.Context) {
	err := recover()
	if err != nil {
		ctx.JSON(GetErrorStatusCode(ErrInternalServer), NewResponse(nil, ErrInternalServer))
		return
	}
}
