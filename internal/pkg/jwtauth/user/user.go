package user

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/winc-link/hummingbird/internal/pkg/jwtauth"
)

func ExtractClaims(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get(jwt.JwtPayloadKey)
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}

func Get(c *gin.Context, key string) interface{} {
	data := ExtractClaims(c)
	if data[key] != nil {
		return data[key]
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " Get 缺少 " + key)
	return nil
}

func GetUserId(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return (data["identity"]).(string)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserId 缺少 identity")
	return ""
}

func GetUserIdStr(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return (data["identity"]).(string)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserIdStr 缺少 identity")
	return ""
}

func GetUserName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["nice"] != nil {
		return (data["nice"]).(string)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetUserName 缺少 nice")
	return ""
}

func GetRoleName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["rolekey"] != nil {
		return (data["rolekey"]).(string)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleName 缺少 rolekey")
	return ""
}

func GetRoleId(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["roleid"] != nil {
		i := (data["roleid"]).(string)
		return i
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetRoleId 缺少 roleid")
	return ""
}

func GetDeptId(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["deptid"] != nil {
		i := (data["deptid"]).(string)
		return i
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetDeptId 缺少 deptid")
	return ""
}

func GetDeptName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["deptkey"] != nil {
		return (data["deptkey"]).(string)
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " [WARING] " + c.Request.Method + " " + c.Request.URL.Path + " GetDeptName 缺少 deptkey")
	return ""
}
