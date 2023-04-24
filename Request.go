package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// This function parses the request URL and returns the parameter
// given by name in param. E.g. /project/{PprojectID} --> "ProjecID" is the URL param
func ParseURLParamInt64(req *http.Request, param string) (int64, error) {
	urlParam := mux.Vars(req)[param]
	value, err := strconv.ParseInt(urlParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Unable to parse URL parameter %s, expected integer, got '%s': %s", param, urlParam, err.Error())
	}

	return value, nil

}
