package auth

import (
	"os"
	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNINGKEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func GetNewToken(id uint, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["email"] = email
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	return tokenString, nil
}
