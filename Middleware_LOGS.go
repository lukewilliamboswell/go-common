package common

import (
	"fmt"
	"net/http"
	"time"
)

func LOGS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		fmt.Println(fmt.Sprintf(
			`%s %s [%s] "%s %s"`,
			req.RemoteAddr,
			req.Host,
			time.Now().Format(time.RFC3339),
			req.Method,
			req.URL,
		))

		next.ServeHTTP(rw, req)
	})
}
