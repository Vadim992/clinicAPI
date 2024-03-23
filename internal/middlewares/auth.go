package middlewares

import (
	"errors"
	"github.com/Vadim992/clinicAPI/internal/config"
	"github.com/Vadim992/clinicAPI/internal/helpers"
	"github.com/Vadim992/clinicAPI/internal/jwtgen"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientTokenSlice := r.Header["Authorization"]

		if clientTokenSlice == nil || len(clientTokenSlice) == 0 {
			logger.InfoLog.Println("client's token is nil or empty string")
			helpers.ClientErr(w, http.StatusUnauthorized)
			return
		}

		clientTokenSliceString := clientTokenSlice[0]
		clientTokenSlice = strings.Split(clientTokenSliceString, "Bearer ")

		if len(clientTokenSlice) != 2 {
			logger.InfoLog.Println("incorrect data in 'Authorization' header request")
			helpers.ClientErr(w, http.StatusUnauthorized)
			return
		}

		clientToken := strings.TrimSpace(clientTokenSlice[len(clientTokenSlice)-1])
		//TODO add checking access token on Blacklist
		claims := jwt.MapClaims{}

		parsedToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AuthCfg.AccessKey), nil
		})

		if err != nil {
			if !errors.Is(err, jwt.ErrTokenExpired) {
				logger.ErrLog.Println("err on parse token with claims:", err)
				helpers.ClientErr(w, http.StatusUnauthorized)
				return
			}

			err = jwtgen.RefreshTokens(w, r, claims)
			if err != nil {
				logger.ErrLog.Println("err on refresh token:", err)
				helpers.ClientErr(w, http.StatusUnauthorized)
			}
			next.ServeHTTP(w, r)

			return
		}

		if !parsedToken.Valid {
			logger.InfoLog.Println("invalid parsed token")
			helpers.ClientErr(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
