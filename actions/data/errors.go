package data

import (
	"net/http"
)

type Data map[string]interface{}

var (
	ErrorNotFound = err(Data{
		"code":    http.StatusNotFound,
		"message": "page not found",
	})
	ErrorUnauthorized = err(Data{
		"code":    http.StatusUnauthorized,
		"message": "authorization required for this action",
	})
	ErrorInternal = err(Data{
		"code":    http.StatusInternalServerError,
		"message": "internal error",
	})
)

func err(data Data) Data {
	return Data{"error": data}
}

func ErrorBadRequest(e error) Data {
	return err(Data{
		"code":    http.StatusBadRequest,
		"message": e.Error(),
	})
}
