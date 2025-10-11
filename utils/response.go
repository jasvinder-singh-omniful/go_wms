// this response format is necessory for consistent interservice
// communication (as specified by go_commons docs)
// https://github.com/omniful/go_commons/tree/master/interservice-client#response-structure
package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
)

type InterSvcResponse struct {
	IsSuccess  bool
	StatusCode int
	Data       json.RawMessage
	Meta       map[string]interface{}
	Error      json.RawMessage
}

func SuccessReponse(c *gin.Context, status http.StatusCode, data interface{}) {
	dataBytes, _ := json.Marshal(data)

	response := InterSvcResponse{
		IsSuccess:  true,
		StatusCode: int(status),
		Data:       json.RawMessage(dataBytes),
		Meta:       make(map[string]interface{}),
		Error:      nil,
	}

	c.JSON(int(status), response)
}

func SendErrorResponse(c *gin.Context, status http.StatusCode, errorMsg string, errors map[string]string) {
    errorData := map[string]interface{}{
        "message": errorMsg,
        "errors":  errors,
    }
    errorBytes, _ := json.Marshal(errorData)
    
    response := InterSvcResponse{
        IsSuccess:  false,
        StatusCode: int(status),
        Data:       nil,
        Meta:       make(map[string]interface{}),
        Error:      json.RawMessage(errorBytes),
    }
    
    c.JSON(int(status), response)
}
