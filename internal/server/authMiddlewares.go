package server

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"personal_budget_app/internal/functionalities"
)


func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
				return
			}
			functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: err.Error()})
			return
		}


		tknStr := c.Value
		claims := &Claims{}

		jwtSecret := os.Getenv("JWT_TOKEN")

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}

			return []byte(jwtSecret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
				return
			}
			functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: err.Error()})
			return
		}
		if !tkn.Valid {
			functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
