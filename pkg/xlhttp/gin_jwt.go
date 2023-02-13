package xlhttp

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/seanlan/xlvein/pkg/xlerror"
	"time"
)

const (
	JWTIdentityKey   = "jwt_user_id"
	RequestTokenHEAD = "X-TOKEN"
)

// JWTBodyMiddleware jwt token 位于请求内
func JWTBodyMiddleware(jwt *JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := Build(c)
		var req struct {
			Token string `form:"token" json:"token" binding:"required"`
		}
		err := r.RequestParser(&req)
		if err != nil {
			r.ctx.Abort()
			return
		}
		jwtUid, err := jwt.ParseToken(req.Token)
		if err != nil {
			r.JsonReturn(xlerror.ErrToken)
			r.ctx.Abort()
			return
		}
		if jwtUid == "" {
			r.JsonReturn(xlerror.ErrToken)
			r.ctx.Abort()
			return
		}
		r.ctx.Set(JWTIdentityKey, jwtUid)
		r.ctx.Next()
	}
}

// JWTHeadMiddleware jwt token 位于请求头
func JWTHeadMiddleware(jwt *JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := Build(c)
		jwtUid, err := jwt.ParseToken(c.GetHeader(RequestTokenHEAD))
		if err != nil {
			r.JsonReturn(xlerror.ErrToken)
			r.ctx.Abort()
			return
		}
		if jwtUid == "" {
			r.JsonReturn(xlerror.ErrToken)
			r.ctx.Abort()
			return
		}
		r.ctx.Set(JWTIdentityKey, jwtUid)
		r.ctx.Next()
	}
}

type JWT struct {
	SigningKey []byte
	Expire     time.Duration
}

func NewJWT(secretKey string, d time.Duration) *JWT {
	return &JWT{
		SigningKey: []byte(secretKey),
		Expire:     d,
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(data string) (string, error) {
	now := time.Now()
	claims := &jwt.StandardClaims{
		Audience:  "",
		IssuedAt:  now.Unix(),
		Issuer:    "TimeToken",
		NotBefore: 0,
		Subject:   data,
	}
	if j.Expire > 0 {
		claims.ExpiresAt = now.Add(j.Expire).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(token string) (string, error) {
	var err error
	var claims jwt.StandardClaims
	_, err = jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		})
	if err != nil {
		return "", err
	}
	err = claims.Valid()
	if err != nil {
		return "", err
	} else {
		return claims.Subject, err
	}
}
