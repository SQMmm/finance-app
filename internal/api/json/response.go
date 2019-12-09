package json

import (
	"encoding/json"
	"net/http"
)

// WriteJsonResponse write generic json response with status code and serializable date
func WriteResponse(resp http.ResponseWriter, statusCode int, data interface{}) error {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	enc := json.NewEncoder(resp)

	return enc.Encode(data)
}
