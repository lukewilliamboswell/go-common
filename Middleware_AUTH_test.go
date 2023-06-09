package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

func TestAUTH(t *testing.T) {
	secret := "my_secret_key"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})
	authMiddleware := AUTH(testHandler, secret)

	// Generate a valid JWT token.
	permissions := []GroupRoles{
		{
			GroupID: 1,
			RoleIDs: []int{1, 2},
		},
	}
	permissionsBytes, _ := json.Marshal(permissions)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"permissions": string(permissionsBytes),
		"exp":         jwt.At(time.Now().Add(time.Hour)),
		"iat":         jwt.At(time.Now()),
	})

	tokenString, _ := token.SignedString([]byte(secret))

	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{"Valid token", tokenString, http.StatusOK},
		{"Invalid token", "invalid_token", http.StatusUnauthorized},
		{"Missing token", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("x-authorization", "Bearer "+tt.token)

			rr := httptest.NewRecorder()
			authMiddleware.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			if tt.expectedCode == http.StatusOK && !strings.Contains(rr.Body.String(), "Hello, world!") {
				t.Error("Expected response body to contain 'Hello, world!'")
			}
		})
	}
}
