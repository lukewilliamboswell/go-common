package common

import (
	"fmt"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Authorization")

		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprint(rw, "")
		} else {
			next.ServeHTTP(rw, req)
		}
	})
}
