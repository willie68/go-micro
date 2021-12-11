package api

import (
	"net/http"
)

// RoleCheck implements a simple middleware handler for adding basic http auth to a route.
func RoleCheck(allowedRoles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			return // TODO deactivated
		})
	}
}
