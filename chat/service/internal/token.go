package internal

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/panjiawan/note/chat/conf"
	"github.com/valyala/fasthttp"
	"time"
)

type CustomClaims struct {
	Uid uint64
	jwt.StandardClaims
}

var expire = 86400

func GenToken(uid uint64) string {
	customClaims := &CustomClaims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expire) * time.Second).Unix(), // 过期时间，必须设置
			Issuer:    "yqyn",
		},
	}

	//采用HMAC SHA256加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	httpConf := conf.GetHandle().GetSysConf()
	tokenString, err := token.SignedString([]byte(httpConf.JwtSecret))
	if err != nil {
		return ""
	}

	return tokenString
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token parse failure:%s", tokenString)
		}
		httpConf := conf.GetHandle().GetSysConf()
		return []byte(httpConf.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func VerifyAuth(ctx *fasthttp.RequestCtx) bool {
	token := ctx.Request.Header.Peek("token")

	if len(token) == 0 {
		return false
	}

	tokenString := string(token)
	if claim, err := ParseToken(tokenString); err == nil {
		ctx.SetUserValue("uid", claim.Uid)
		return true
	}

	return false
}

func GetUid(ctx *fasthttp.RequestCtx) (uint64, error) {
	if uid := ctx.UserValue("uid"); uid != nil {
		// 如果存在uid数据
		return uid.(uint64), nil
	} else {
		// 不存在则试图从token中获取
		token := ctx.Request.Header.Peek("token")
		if len(token) == 0 {
			return 0, errors.New("uid empty")
		}
		tokenString := string(token)
		if claim, err := ParseToken(tokenString); err == nil {
			return claim.Uid, nil
		}
	}

	return 0, errors.New("uid empty")
}
