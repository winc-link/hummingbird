package jwt

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"github.com/winc-link/hummingbird/internal/pkg/middleware"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	JwtIssuer = "hummingbird"
	tokenKey  = "x-token"
)

var (
	RefreshKey = "afe60362-r3f2-5e07-4d9f-g35e32650af3"
	OpenAPIKey string
	JwtSignKey string
)

func SetOpenAPIKey(key string) {
	OpenAPIKey = key
}

func SetJwtSignKey(key string) {
	JwtSignKey = key
}

func GenerateJwtSignKey() string {
	return uuid.New().String()
}

//JWT 令牌认证中间件
func JWTAuth(CloseAuthToken bool) gin.HandlerFunc {
	if CloseAuthToken {
		return func(c *gin.Context) {}
	}
	return func(c *gin.Context) {
		token := c.Request.Header.Get(tokenKey)
		vars := c.Request.URL.Query()
		if v, ok := vars[tokenKey]; ok {
			token = v[0]
		}
		if token == "" {
			httphelper.RenderFailNoLog(c, TokenInvalid, c.Writer)
			c.Abort()
			return
		}
		j := NewJWT(JwtSignKey)
		claims, err := j.ParseToken(token)
		if err != nil {
			httphelper.RenderFailNoLog(c, TokenExpired, c.Writer)
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt {
			httphelper.RenderFailNoLog(c, TokenExpired, c.Writer)
			c.Abort()
			return
		}
		c.Set(constants.JwtParsedInfo, claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errort.NewCommonErr(errort.DefaultTokenPermission, fmt.Errorf("Token expired, please log in again"))
	TokenNotValidYet = errort.NewCommonErr(errort.DefaultTokenPermission, fmt.Errorf("Token expired, please log in again"))
	TokenMalformed   = errort.NewCommonErr(errort.DefaultTokenPermission, fmt.Errorf("Token expired, please log in again"))
	TokenInvalid     = errort.NewCommonErr(errort.DefaultTokenPermission, fmt.Errorf("Token expired, please log in again"))
)

func NewJWT(key string) *JWT {
	return &JWT{
		[]byte(key),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims middleware.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*middleware.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*middleware.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
