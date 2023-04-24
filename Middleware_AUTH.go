package common

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func AUTH(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			next.ServeHTTP(rw, req)
			return
		}

		// Token will be in either header or as a query param
		var tokenString string
		authHeader := req.Header.Get("x-authorization")
		if authHeader == "" {
			query := req.URL.Query()
			if len(query["x-authorization"]) > 0 {
				tokenString = query["x-authorization"][0]
			}
		} else {
			tokenString = authHeader[7:]
		}

		if len(tokenString) <= 0 {
			BadRequestResponse(fmt.Errorf("missing token in x-authorization header")).ServeHTTP(rw, req)
			return
		}

		// parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			UnauthorisedResponse(fmt.Errorf("expected a valid token in x-authorization header")).ServeHTTP(rw, req)
			return
		}

		// parse permissions from token claim
		permissions := parseGroupRoles(token.Claims["permissions"])
		if permissions == nil {
			UnauthorisedResponse(fmt.Errorf("expected a valid token in x-authorization header")).ServeHTTP(rw, req)
			return
		}

		// add permissions to context
		SetPermissions(req, permissions)

		// continue with middleware chain
		next.ServeHTTP(rw, req)
		return
	})
}