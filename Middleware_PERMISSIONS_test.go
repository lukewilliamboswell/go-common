package common

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMustHaveGroupRole(t *testing.T) {
	// Test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access granted"))
	})

	tests := []struct {
		name         string
		permissions  []GroupRoles
		groupID      int
		roleID       int
		expectedCode int
	}{
		{
			name: "Access granted",
			permissions: []GroupRoles{
				{GroupID: 1, RoleIDs: []int{1, 2}},
			},
			groupID:      1,
			roleID:       1,
			expectedCode: http.StatusOK,
		},
		{
			name: "Access denied",
			permissions: []GroupRoles{
				{GroupID: 1, RoleIDs: []int{1, 2}},
			},
			groupID:      1,
			roleID:       3,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			SetPermissions(req, tt.permissions)

			rr := httptest.NewRecorder()
			middleware := MustHaveGroupRole(testHandler, tt.groupID, tt.roleID)
			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}
		})
	}
}
