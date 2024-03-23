package middlewares

import (
	"fmt"
	"github.com/Vadim992/clinicAPI/internal/config"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/helpers"
	"github.com/Vadim992/clinicAPI/internal/helpers/dbhelpers"
	"github.com/Vadim992/clinicAPI/internal/jwtgen"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func CheckPermissions(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessTokenString := jwtgen.GetAccessTokenFromRequest(w, r)

		if accessTokenString == "" {
			return
		}

		claims := jwt.MapClaims{}

		parsedToken, err := jwt.ParseWithClaims(accessTokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AuthCfg.AccessKey), nil
		})

		if err != nil {
			logger.ErrLog.Println("err on parse token with claims:", err)

			return
		}

		if !parsedToken.Valid {
			logger.InfoLog.Println("invalid parsed token")

			return
		}

		id, err := jwtgen.GetIntFromTokenField(claims, "sub")

		if err != nil {
			logger.ErrLog.Println("cant get id from claims:", err)

			return
		}

		roleId, err := jwtgen.GetIntFromTokenField(claims, "role_id")

		if err != nil {
			logger.ErrLog.Println("cant get id from claims:", err)

			return
		}

		idReq, err := dbhelpers.ConvertIdFromStrToInt(r)

		if err != nil {
			helpers.ServeErr(w, err)
			return
		}

		if idReq != id {
			helpers.ClientErr(w, http.StatusBadRequest)
			return
		}

		isPermitted, err := hasPermission(roleId, "/patients/id", r.Method)

		if err != nil {
			helpers.ServeErr(w, err)
			return
		}

		if !isPermitted {
			helpers.ClientErr(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func hasPermission(roleId int, path, method string) (bool, error) {
	permissions, err := postgres.DataBase.GetPatientIdPermissions(roleId, path)
	fmt.Println(permissions)

	if err != nil {
		return false, err
	}

	for _, val := range permissions {
		if val == method {
			return true, nil
		}
	}

	return false, nil
}
