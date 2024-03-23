package jwtgen

import (
	"github.com/Vadim992/clinicAPI/internal/config"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func generateAccessToken(role, id int) (string, error) {
	expirationTime := time.Now().Add(20 * time.Second).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Header["kid"] = "access_token"

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTime
	claims["role_id"] = role
	claims["sub"] = id

	tokenString, err := token.SignedString([]byte(config.AuthCfg.AccessKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GeneratePairsToken(role, id int) (string, string, error) {
	accessTokenString, err := generateAccessToken(role, id)

	if err != nil {
		return "", "", err
	}

	expirationTime := time.Now().Add(40 * time.Second).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Header["kid"] = "refresh_token"

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTime
	claims["role_id"] = role
	claims["sub"] = id

	refreshTokenString, err := token.SignedString([]byte(config.AuthCfg.RefreshKey))

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func RefreshTokens(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) error {
	roleId, err := GetIntFromTokenField(claims, "role_id")
	if err != nil {
		return err
	}

	id, err := GetIntFromTokenField(claims, "sub")
	if err != nil {
		return err
	}

	switch {
	case roleId == 1 || roleId == 3:
		if err := refreshPatientTokens(w, r, id, roleId); err != nil {
			return err
		}
	case roleId == 2:
		//TODO for doctor
	default:
		return jwt.ErrTokenInvalidClaims
	}

	return nil
}

func refreshPatientTokens(w http.ResponseWriter, r *http.Request, id, roleId int) error {
	err := validateToken(id)

	if err != nil {
		return err
	}

	_, refreshTokenString, err := GeneratePairsToken(id, roleId)

	if err != nil {
		return err
	}

	err = postgres.DataBase.UpdatePatientRefreshToken(id, refreshTokenString)

	if err != nil {
		return err
	}

	//jwtokens := dto.NewJWTokens(accessTokenString, refreshTokenString)
	//tokens, err := jwtokens.EncodeToJSON()
	//
	//if err != nil {
	//	return err
	//}

	//reqHeader := fmt.Sprintf("Bearer %s", accessTokenString)
	//r.Header.Set("Authorization", reqHeader)
	//w.Write(tokens)

	return nil
}

func validateToken(id int) error {
	refreshTokenString, err := postgres.DataBase.CheckPatientRefreshToken(id)

	if err != nil {
		return nil
	}

	parsedToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AuthCfg.RefreshKey), nil
	})

	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return jwt.ErrTokenInvalidId
	}

	return nil
}

func GetIntFromTokenField(claims jwt.MapClaims, field string) (int, error) {
	numFloat, ok := claims[field].(float64)

	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	num := int(numFloat)

	return num, nil
}

func GetAccessTokenFromRequest(w http.ResponseWriter, r *http.Request) string {
	clientTokenSlice := r.Header["Authorization"]

	if clientTokenSlice == nil || len(clientTokenSlice) == 0 {
		logger.InfoLog.Println("client's token is nil or empty string")
		w.WriteHeader(http.StatusUnauthorized)
		return ""
	}

	clientTokenSliceString := clientTokenSlice[0]
	clientTokenSlice = strings.Split(clientTokenSliceString, "Bearer ")

	if len(clientTokenSlice) != 2 {
		logger.InfoLog.Println("incorrect data in 'Authorization' header request")
		w.WriteHeader(http.StatusUnauthorized)
		return ""
	}

	clientToken := strings.TrimSpace(clientTokenSlice[len(clientTokenSlice)-1])

	return clientToken
}
