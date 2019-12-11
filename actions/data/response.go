package data

import (
	"github.com/sqmmm/finance-app/internal/api/json"
	"github.com/sqmmm/finance-app/internal/logger"
	"net/http"
)

func Result(result interface{}) Data {
	return Data{"result": result}
}

func WriteResponse(log logger.Logger, resp http.ResponseWriter, statusCode int, data interface{}) {
	err := json.WriteResponse(resp, statusCode, data)
	if err != nil {
		log.Errorf("failed to write response: %s", err)
	}
}
