package common

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

type key int

const PermissionsKey key = 123234

func parseGroupRoles(v interface{}) []GroupRoles {

	if v == nil {
		return nil
	}

	// convert to json
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	// parse back
	permissions := make([]GroupRoles, 0)
	err = json.Unmarshal(data, &permissions)
	if err != nil {
		return nil
	}

	return permissions

}

type GroupRoles struct {
	GroupID int   `json:"groupID"`
	RoleIDs []int `json:"roleIDs"`
}

func SetPermissions(req *http.Request, permissions []GroupRoles) {
	context.Set(req, PermissionsKey, permissions)
}

func GetPermissions(req *http.Request) []GroupRoles {

	if rv := context.Get(req, PermissionsKey); rv != nil {
		return rv.([]GroupRoles)
	}

	return nil

}

var ErrIncorrectPermissions error = errors.New("incorrect permissions, missing group or role to perform this action")

func MustHaveGroupRole(next http.Handler, groupID, roleID int) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// get permmissions from request context
		// must have had auth middle ware called previously
		permissions := GetPermissions(req)
		if permissions == nil {
			// expected a permissions object
			log.Println("WARNING: expected permissions to be set in context")
			next.ServeHTTP(rw, req)
			return
		}

		for i, _ := range permissions {

			// ignore other groupIDs
			if permissions[i].GroupID != groupID {
				continue
			}

			for j, _ := range permissions[i].RoleIDs {

				if permissions[i].RoleIDs[j] != roleID {
					continue
				}

				// continue with middleware chain
				next.ServeHTTP(rw, req)
				return

			}
		}

		// couldnt find the role in any groups...
		UnauthorisedResponse(ErrIncorrectPermissions).ServeHTTP(rw, req)

	})
}
