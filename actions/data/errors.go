package data

import (
	"net/http"
)

type Data map[string]interface{}

var (
	ErrorNotFound = err(Data{
		"code":    http.StatusBadRequest,
		"message": "not found",
	})
)

func err(data Data) Data {
	return Data{"error": data}
}
