package token

import (
	"time"

	"github.com/microservice/api_gateway/pkg/logger"

	"github.com/dgrijalva/jwt-go"
)

type JWTHandler struct {
	Sub       string
	Iss       string
	Exp       string
	Iat       string
	Aud       []string
	Role      string
	SigninKEY string
	Log       logger.Logger
	Token     string
}

func (jwtHandler *JWTHandler) GenerateAuthJWT() (acces, refresh string, err error) {
	var (
		accesToken   *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
	)

	accesToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accesToken.Claims.(jwt.MapClaims)
	claims["iss"] = jwtHandler.Iss
	claims["sub"] = jwtHandler.Sub
	claims["exp"] = time.Now().Add(time.Hour * 500).Unix()
	claims["iat"] = time.Now().Unix()
	claims["role"] = jwtHandler.Role
	claims["aud"] = jwtHandler.Aud

	acces, err = accesToken.SignedString([]byte(jwtHandler.SigninKEY))
	if err != nil {
		jwtHandler.Log.Error("error generating acces token", logger.Error(err))
		return
	}

	refresh, err = refreshToken.SignedString([]byte(jwtHandler.SigninKEY))
	if err != nil {
		jwtHandler.Log.Error("error generating refresh token", logger.Error(err))
		return
	}

	return acces, refresh, nil
}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SigninKEY), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		jwtHandler.Log.Error("invalid jwt token")
	}

	return claims, nil
}
