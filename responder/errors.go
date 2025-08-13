package responder

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error        any    `json:"error"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    int    `json:"error_code"`
}

func SendError(w http.ResponseWriter, errMessage string, status int, err ...error) {
	errResp := ErrorResponse{
		Error:        nil,
		ErrorMessage: errMessage,
		ErrorCode:    1000,
	}
	if len(err) > 0 && err[0] != nil {
		errResp.Error = err[0].Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}

func SendErrorWithParams(w http.ResponseWriter, err string, status int, errorCode *int, errorMessage *string) {
	errResp := ErrorResponse{
		Error:        err,
		ErrorMessage: "",
		ErrorCode:    1000,
	}

	if errorCode != nil && *errorCode > 0 {
		errResp.ErrorCode = *errorCode
	}

	if errorMessage != nil && len(*errorMessage) > 0 {
		errResp.ErrorMessage = *errorMessage
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}

func BadBody(w http.ResponseWriter, err error) {
	SendError(w, "bad body", http.StatusBadRequest, err)
}

func ErrMissingBodyRequirement(w http.ResponseWriter, field string) {
	SendError(w, "missing required body field: "+field, http.StatusBadRequest)
}

func ErrInvalidBodyField(w http.ResponseWriter, field string, err error) {
	SendError(w, "invalid body field: "+field, http.StatusBadRequest, err)
}

func ErrConflict(w http.ResponseWriter, err error) {
	SendError(w, "conflict", http.StatusConflict, err)
}

func ErrInternal(w http.ResponseWriter, err error, message string) {
	SendError(w, message, http.StatusInternalServerError, err)
}

func ParamError(w http.ResponseWriter, field string) {
	SendError(w, "missing required param: "+field, http.StatusBadRequest)
}

func ErrRequiredKey(w http.ResponseWriter, key string) {
	SendError(w, "missing required key: "+key, http.StatusBadRequest)
}
