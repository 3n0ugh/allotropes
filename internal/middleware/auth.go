package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/3n0ugh/allotropes/internal/token"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim := token.Claims{}
		err := claim.ParseToken(r.Header.Get("authorization"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			res, _ := json.Marshal(errors.NewUnAuthorizedError("unauthorized", err.Error()))
			w.Write(res)
			return
		}

		ctx := context.WithValue(r.Context(), "email", claim.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
