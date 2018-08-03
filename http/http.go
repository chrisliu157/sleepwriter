package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type ErrorType string

var Error = struct {
	Validation     ErrorType
	NotFound       ErrorType
	InternalServer ErrorType
	Duplicate      ErrorType
	Unauthorized   ErrorType
}{
	"validation_error",
	"not_found",
	"internal_server_error",
	"already_exists",
	"unauthorized",
}

// Parse HTTP Body
func ParseBody(w http.ResponseWriter, request io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(request)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// ResponseWriter is a interface, we should pass by value because it internally contains a pointer to the actual writer
func RenderError(w http.ResponseWriter, code int, error_code ErrorType, error_msg string) {
	format := make(map[string]interface{})
	format["error_code"] = error_code
	format["error"] = error_msg

	response, _ := json.Marshal(format)

	w.WriteHeader(code)
	w.Write(response)
}

func RenderAsync(w http.ResponseWriter, job_id string) {
	job := make(map[string]string)
	job["location"] = fmt.Sprintf("/api/v1/jobs/%s", job_id)

	response, _ := json.Marshal(job)

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
