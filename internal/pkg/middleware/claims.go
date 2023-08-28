package middleware

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID       uint
	Username string
	jwt.StandardClaims
}
