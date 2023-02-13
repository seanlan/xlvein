package veinsdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/parnurzeal/gorequest"
	"github.com/seanlan/goutils/xljson"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func MakeTimedToken(data, secret string, expire int64) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:  "TimeToken",
		Subject: data,
	}
	if expire > 0 {
		claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(expire)).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func MapToUrlencoded(m map[string]interface{}, secretKey string) string {
	var keys []string
	var _source []string
	for k := range m {
		keys = append(keys, k)
	}
	//字符串排序
	sort.Strings(keys)
	for _, k := range keys {
		_source = append(_source, fmt.Sprintf("%s=%v", k, m[k]))
	}
	//map URL加入密钥拼接
	_source = append(_source, fmt.Sprintf("%s=%s", "key", secretKey))
	sourceStr := strings.Join(_source, "&")
	//MD5加密
	h := md5.New()
	h.Write([]byte(sourceStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

type SDK struct {
	Gateway   string
	AppID     string
	AppSecret string
}

func New(gate, appid, secret string) *SDK {
	return &SDK{
		Gateway:   gate,
		AppID:     appid,
		AppSecret: secret,
	}
}

// GetSign 参数签名
func (client *SDK) GetSign(jsonObject map[string]interface{}) string {
	return MapToUrlencoded(jsonObject, client.AppSecret)
}

// GetNonce 获取随机字符串
func (client *SDK) GetNonce() string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ProduceIMToken 获取用户Token
func (client *SDK) ProduceIMToken(userTag string) (string, error) {
	return MakeTimedToken(userTag, client.AppSecret, 3600*24*30)
}

// Request 基础请求接口
func (client *SDK) Request(path string, jsonObject map[string]interface{}) (string, error) {
	request := gorequest.New()
	_, body, errs := request.Post(path).
		Type("form").
		Send(jsonObject).End()
	var err error
	if len(errs) > 0 {
		err = errs[0]
	} else {
		err = nil
	}
	return body, err
}

// ApiRequest API封装请求
func (client *SDK) ApiRequest(method string, jsonObject map[string]interface{}) (resp xljson.JsonObject, err error) {
	//API接口请求
	var buf bytes.Buffer
	buf.WriteString(client.Gateway)
	buf.WriteString(method)
	apiUrl := buf.String()
	sign := MapToUrlencoded(jsonObject, client.AppSecret)
	jsonObject["sign"] = sign
	body, err := client.Request(apiUrl, jsonObject)
	if err != nil {
		return
	}
	resp = xljson.JsonObject{Buff: []byte(body)}
	return
}

func (client *SDK) PushMessage(sendTo string, message map[string]interface{}) (resp xljson.JsonObject, err error) {
	var msg []byte
	msg, err = json.Marshal(message)
	if err != nil {
		return
	}
	return client.ApiRequest("/api/v1/im/push", map[string]interface{}{
		"app_key": client.AppID,
		"send_to": sendTo,
		"message": string(msg),
		"nonce":   client.GetNonce(),
	})
}
