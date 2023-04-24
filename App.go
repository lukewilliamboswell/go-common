package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	ErrRequiredField     error = errors.New("*Required field")
	ErrUnsupportedValue  error = errors.New("unsupported Value")
	ErrRecordDoesntExist error = errors.New("record does not exist")
	ErrNotImplemented    error = errors.New("Endpoint not implemented")
	ErrBadJSONBody       error = errors.New("unable to parse body, exptected json")
)

type FieldError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type FieldErrors []FieldError

func (e *FieldErrors) Append(field, message string) {
	*e = append(*e, FieldError{
		Field:   field,
		Message: message,
	})
}

type ResponseType int

const (
	RESPOND_OK ResponseType = iota
	RESPOND_BAD_REQUEST
	RESPOND_INTERNAL_ERROR
	RESPOND_USER_ERROR
	RESPOND_DELETED
	RESPOND_UPDATED
	RESPOND_PAGED
	RESPOND_FILE_DOWNLOAD
	RESPOND_UNAUTHORISED
)

type SubResponseDeleted struct {
	Count int64 `json:"count"`
}

type SubResponseUpdated struct {
	Count int64 `json:"count"`
}

type SubResponsePaged struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"perPage"`
	Total   int64 `json:"totalCount"`
}

type Response struct {
	Status   bool                `json:"success"`
	Type     ResponseType        `json:"-"`
	Result   interface{}         `json:"response,omitempty"`
	Error    string              `json:"error,omitempty"`
	Errors   FieldErrors         `json:"errors,omitempty"`
	Updated  *SubResponseUpdated `json:"updated,omitempty"`
	Deleted  *SubResponseDeleted `json:"deleted,omitempty"`
	Paged    *SubResponsePaged   `json:"paged,omitempty"`
	File     *os.File            `json:"-"`
	FileName string              `json:"-"`
}

func BadRequestResponse(err error) Response {
	return Response{
		Type:  RESPOND_BAD_REQUEST,
		Error: err.Error(),
	}
}

func UnauthorisedResponse(err error) Response {
	return Response{
		Type:  RESPOND_UNAUTHORISED,
		Error: err.Error(),
	}
}

func BadRequestErrorsResponse(errors FieldErrors) Response {
	return Response{
		Type:   RESPOND_BAD_REQUEST,
		Errors: errors,
	}
}

func UserErrorResponse(err error) Response {
	return Response{
		Type:  RESPOND_USER_ERROR,
		Error: err.Error(),
	}
}

func FileDownloadResponse(file *os.File, filename string) Response {
	return Response{
		Type:     RESPOND_FILE_DOWNLOAD,
		File:     file,
		FileName: filename,
	}
}

func PagedResponse(result interface{}, page, per_page, total int64) Response {
	return Response{
		Type:   RESPOND_PAGED,
		Result: result,
		Paged: &SubResponsePaged{
			Page:    page,
			PerPage: per_page,
			Total:   total,
		},
	}
}

func DeletedResponse(count int64) Response {
	return Response{
		Type: RESPOND_DELETED,
		Deleted: &SubResponseDeleted{
			Count: count,
		},
	}
}

func UpdatedResponse(count int64) Response {
	return Response{
		Type: RESPOND_UPDATED,
		Updated: &SubResponseUpdated{
			Count: count,
		},
	}
}

func InternalErrorResponse(err error) Response {
	return Response{
		Type:  RESPOND_INTERNAL_ERROR,
		Error: err.Error(),
	}
}

func SuccessResult(result interface{}) Response {
	return Response{
		Type:   RESPOND_OK,
		Result: result,
	}
}

func (r Response) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// HEADERS
	switch r.Type {
	default:
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		r.Error = "ERROR: unrecognised response type"
	case RESPOND_BAD_REQUEST:
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
	case RESPOND_INTERNAL_ERROR:
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
	case RESPOND_USER_ERROR, RESPOND_OK, RESPOND_DELETED, RESPOND_UPDATED, RESPOND_PAGED:
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	case RESPOND_UNAUTHORISED:
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
	case RESPOND_FILE_DOWNLOAD:
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", r.FileName))
		rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))
		rw.WriteHeader(http.StatusOK)
	}

	// STATUS (request succeded or failed from the client's perspective)
	switch r.Type {
	default:
		r.Status = false
	case RESPOND_OK, RESPOND_DELETED, RESPOND_UPDATED, RESPOND_PAGED, RESPOND_FILE_DOWNLOAD:
		r.Status = true
	}

	// SERIALISE RESPONSE
	switch r.Type {
	default:

		err := json.NewEncoder(rw).Encode(r)
		if err != nil {
			log.Println("ERROR: Unable to marshal JSON error", err)
		}

	case RESPOND_FILE_DOWNLOAD:
		_, err := io.Copy(rw, r.File)
		if err != nil {
			log.Println("ERROR: Unable to write file response", err)
		}

		r.File.Close()
	}
}

type AppHandler func(*http.Request) http.Handler

func (fn AppHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fn(req).ServeHTTP(rw, req)
}
