package userapp

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	pkgcontainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/middleware"
	"time"

	//"gitlab.com/tedge/edgex/internal/pkg/container"
	//resourceContainer "gitlab.com/tedge/edgex/internal/tedge/resource/container"
	//
	//"gitlab.com/tedge/edgex/internal/pkg/di"
	//"gitlab.com/tedge/edgex/internal/pkg/errort"
	//"gitlab.com/tedge/edgex/internal/pkg/logger"
	//
	jwt2 "github.com/winc-link/hummingbird/internal/tools/jwt"
	//
	//"github.com/dgrijalva/jwt-go"
	//"gitlab.com/tedge/edgex/internal/dtos"
	//"gitlab.com/tedge/edgex/internal/models"
	//"gitlab.com/tedge/edgex/internal/pkg/middleware"
	//"gitlab.com/tedge/edgex/internal/tedge/resource/interfaces"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultUserName = "admin"
	DefaultLang     = "en"
)

var _ interfaces.UserItf = new(userApp)

type userApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func New(dic *di.Container) *userApp {
	return &userApp{
		dic:      dic,
		lc:       pkgcontainer.LoggingClientFrom(dic.Get),
		dbClient: container.DBClientFrom(dic.Get),
	}
}

//UserLogin 用户登录
func (uapp *userApp) UserLogin(ctx context.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	// 从数据库用户信息
	user, edgeXErr := uapp.dbClient.GetUserByUserName(req.Username)
	if edgeXErr != nil {
		return res, errort.NewCommonEdgeX(errort.AppPasswordError, "", edgeXErr)
	}

	// 校验密码
	cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if cErr != nil {
		err = errort.NewCommonErr(errort.AppPasswordError, cErr)
		return
	}

	j := jwt2.NewJWT(jwt2.JwtSignKey)
	claims := middleware.CustomClaims{
		ID:       1,
		Username: req.Username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*3, // 过期时间 7天
			Issuer:    jwt2.JwtIssuer,                 // 签名的发行者
		},
	}

	token, jwtErr := j.CreateToken(claims)
	if jwtErr != nil {
		err = jwtErr
		return
	}
	lang := user.Lang
	if lang == "" {
		lang = DefaultLang
	}
	res = dtos.LoginResponse{
		User: dtos.UserResponse{
			Username: user.Username,
			Lang:     lang,
		},
		ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		Token:     token,
	}
	return
}

//InitInfo 查询用户信息
func (uapp *userApp) InitInfo() (res dtos.InitInfoResponse, err error) {
	// 从数据库用户信息
	_, edgeXErr := uapp.dbClient.GetUserByUserName(DefaultUserName)
	if edgeXErr != nil {
		//if errort.NewCommonEdgeXWrapper(edgeXErr).Code() ==  {
		if errort.Is(errort.DefaultResourcesNotFound, edgeXErr) {
			res.IsInit = false
			return
		}
		return
	}
	res.IsInit = true
	return
}

// InitPassword 初始化密码
func (uapp *userApp) InitPassword(ctx context.Context, req dtos.InitPasswordRequest) error {
	lc := uapp.lc

	// 从数据库用户信息
	_, edgeXErr := uapp.dbClient.GetUserByUserName(DefaultUserName)
	if edgeXErr == nil {
		return errort.NewCommonErr(errort.AppSystemInitialized, edgeXErr)
	}

	// 生成新密码并存储
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := models.User{
		Username:   DefaultUserName,
		Password:   string(newPasswordHash),
		Lang:       DefaultLang,
		OpenAPIKey: jwt2.GenerateJwtSignKey(),
		GatewayKey: jwt2.GenerateJwtSignKey(),
	}
	jwt2.SetOpenAPIKey(newUser.OpenAPIKey)
	jwt2.SetJwtSignKey(newUser.GatewayKey)
	//db操作存储
	_, edgeXErr = uapp.dbClient.AddUser(newUser)
	if edgeXErr != nil {
		lc.Errorf("add user error %v", edgeXErr)
		return edgeXErr
	}
	return nil
}

// UpdateUserPassword 修改密码
func (uapp *userApp) UpdateUserPassword(ctx context.Context, username string, req dtos.UpdatePasswordRequest) error {
	// 从数据库用户信息
	user, edgeXErr := uapp.dbClient.GetUserByUserName(username)
	if edgeXErr != nil {
		return edgeXErr
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return err
	}

	// 生成新密码并存储
	newPasswordHash, gErr := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if gErr != nil {
		return err
	}
	user.Password = string(newPasswordHash)
	//db操作存储
	edgeXErr = uapp.dbClient.UpdateUser(user)
	if edgeXErr != nil {
		return edgeXErr
	}
	return nil
}

//OpenApiUserLogin openapi用户登录
func (uapp *userApp) OpenApiUserLogin(ctx context.Context, req dtos.LoginRequest) (res *dtos.TokenDetail, err error) {
	// 从数据库用户信息
	user, edgeXErr := uapp.dbClient.GetUserByUserName(req.Username)
	if edgeXErr != nil {
		return res, errort.NewCommonErr(errort.AppPasswordError, edgeXErr)
	}

	// 校验密码
	cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if cErr != nil {
		err = errort.NewCommonErr(errort.AppPasswordError, cErr)
		return
	}
	td, err := uapp.CreateTokenDetail(req.Username)
	if err != nil {
		return
	}
	return td, nil
}

// CreateTokenDetail 根据用户名生成 token
func (uapp *userApp) CreateTokenDetail(userName string) (*dtos.TokenDetail, error) {
	td := &dtos.TokenDetail{
		AccessId:  "accessId",
		RefreshId: "refreshId",
		AtExpires: time.Now().Add(time.Minute * 120).Unix(),   //两小时
		RtExpires: time.Now().Add(time.Hour * 24 * 14).Unix(), //两星期
	}
	var (
		userID uint = 1
		err    error
	)
	td.AccessToken, err = uapp.createToken(userID, userName, td.AtExpires, jwt2.JwtSignKey)
	if err != nil {
		return nil, err
	}
	td.RefreshToken, err = uapp.createToken(userID, userName, td.RtExpires, jwt2.RefreshKey)
	if err != nil {
		return nil, err
	}
	return td, nil
}

// CreateToken 生成 token
func (uapp *userApp) createToken(useId uint, userName string, expire int64, signKey string) (string, error) {
	j := jwt2.NewJWT(signKey)
	claims := middleware.CustomClaims{
		ID:       useId,
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), // 签名生效时间
			ExpiresAt: expire,
			Issuer:    jwt2.JwtIssuer, // 签名的发行者
		},
	}

	token, jwtErr := j.CreateToken(claims)
	if jwtErr != nil {
		err := jwtErr
		return "", err
	}
	return token, nil
}

func (uapp *userApp) InitJwtKey() {
	user, err := uapp.dbClient.GetUserByUserName(DefaultUserName)
	if err != nil {
		return
	}
	if user.GatewayKey == "" {
		user.GatewayKey = jwt2.GenerateJwtSignKey()
	}
	if user.OpenAPIKey == "" {
		user.OpenAPIKey = jwt2.GenerateJwtSignKey()
	}
	if err = uapp.dbClient.UpdateUser(user); err != nil {
		panic(err)
	}
	jwt2.SetOpenAPIKey(user.OpenAPIKey)
	jwt2.SetJwtSignKey(user.GatewayKey)
}
