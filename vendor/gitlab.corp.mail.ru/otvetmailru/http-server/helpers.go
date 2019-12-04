package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// VerboseLevel describes errors detalization level
type VerboseLevel int64

// Logs detalization level
const (
	SimpleErrors VerboseLevel = iota
	DetailedErrors
	StackTrace
)

var verboseLevel VerboseLevel

// SetVerboseLevel changes errors detalization lebel
func SetVerboseLevel(level VerboseLevel) {
	atomic.StoreInt64((*int64)(&verboseLevel), int64(level))
}

// GetVerboseLevel returns currently set verbose level
func GetVerboseLevel() VerboseLevel {
	return VerboseLevel(atomic.LoadInt64((*int64)(&verboseLevel)))
}

var debugMode bool

// EnableDebugMode enables debug mode,
// in debug mode errors sends in response to client
func EnableDebugMode() {
	debugMode = true
}

var errorLogger func(message string)

func RegisterErrorLogger(loggerFunc func(message string)) {
	errorLogger = loggerFunc
}

func writeLog(message string) {
	if errorLogger == nil {
		return
	}

	errorLogger(message)
}

// Data default struct for response (it is easier to write Data{"field": value} than whole map)
type Data map[string]interface{}

// CtxParamInt extract parameter from request context using key
func CtxParamInt(request *http.Request, key string) int64 {
	var param int64
	var err error

	switch val := request.Context().Value(key).(type) {
	case string:
		param, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
	case int, int64, int32:
		param = val.(int64)
	default:
		param = 0
	}

	return param
}

// CtxVarInt extract query argument from request context using key
func CtxVarInt(request *http.Request, key string) int64 {
	ov := mux.Vars(request)[key]
	val, err := strconv.ParseInt(ov, 10, 64)
	if err != nil {
		val = 0
	}
	return val
}

// CtxVarString extract query argument from request context using key
func CtxVarString(request *http.Request, key string) string {
	return mux.Vars(request)[key]
}

func QueryParamInt(request *http.Request, key string) int64 {
	val := request.URL.Query().Get(key)
	res, _ := strconv.ParseInt(val, 10, 64)

	return res
}

func QueryParamString(request *http.Request, key string) string {
	return request.URL.Query().Get(key)
}

// QueryParamTime tries to parses Time value from request parameter.
func QueryParamTime(request *http.Request, key, layout string) (time.Time, error) {
	ts := request.URL.Query().Get(key)

	return time.Parse(layout, ts)
}

// CurrentRoute returns current matched route
func CurrentRoute(request *http.Request) *mux.Route {
	return mux.CurrentRoute(request)
}

// WriteJsonResponse write generic json response with status code and serializable date
func WriteJsonResponse(resp http.ResponseWriter, statusCode int, data interface{}) error {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	enc := json.NewEncoder(resp)

	return enc.Encode(data)
}

// WriteJsonNotFound writes 404 Not Found response to client
func WriteJsonNotFound(resp http.ResponseWriter) error {
	return WriteJsonResponse(resp, 404, Data{"error": "Not found"})
}

// WriteJsonServerError writes 500 IT IS BROKEN PANIC AAAAAAA!!1 response to client
func WriteJsonServerError(resp http.ResponseWriter, err error) error {
	var errMessage string
	if verboseLevel == SimpleErrors {
		errMessage = fmt.Sprintf("%v", err)
	} else {
		errMessage = fmt.Sprintf("%+v", err)
		errMessage = strings.Replace(errMessage, "\t", "", -1)
	}

	writeLog(errMessage)

	var data Data
	if debugMode {
		lines := strings.Split(errMessage, "\n")
		if len(lines) > 1 {
			data = Data{"error": lines[0], "lines": lines[1:]}
		} else {
			data = Data{"error": errMessage}
		}
	} else {
		data = Data{"error": "server error"}
	}

	return WriteJsonResponse(resp, 500, data)
}

func ReadJsonData(request *http.Request, data interface{}) error {
	dec := json.NewDecoder(request.Body)
	return dec.Decode(data)
}
