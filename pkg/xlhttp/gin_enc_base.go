/**
	gin 数据加密中间件
支持request请求加密、response响应数据加密
支持加密方式有
 AES 对称加密
也可以自己实现 Encryptor 接口来自定义加密方式

request请求加密方式，将请求的json加密后，以表单的形式提交
enc=请求json加密信息

response响应数据加密，将返回的数据加密后返回

*/

package xlhttp

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strings"
)

type EncRequest struct {
	Enc string `json:"enc" form:"enc" binding:"required"`
}

type Encryptor interface {
	Decrypt(source string) (string, error) //加密
	Encrypt(dec string) (string, error)    //解密
}

func NewEncryptRequestMiddleware(encryptor Encryptor, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			data string
			err  error
			req  EncRequest
		)
		contentType := c.ContentType()
		method := c.Request.Method
		switch method {
		case http.MethodPost:
			switch {
			case strings.Contains(contentType, gin.MIMEJSON):
				err = c.ShouldBindBodyWith(&req, binding.JSON)
				break
			case strings.Contains(contentType, gin.MIMEPOSTForm),
				strings.Contains(contentType, gin.MIMEPOSTForm):
				err = c.ShouldBindWith(&req, binding.Form)
				break
			default:
				err = c.ShouldBind(&req)
			}
			break
		case http.MethodGet:
			err = c.ShouldBindQuery(&req)
		}
		if err != nil && !debug {
			c.AbortWithStatus(500)
			return
		}
		enc := req.Enc
		if len(enc) > 0 || !debug {
			data, err = encryptor.Decrypt(enc)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			c.Request.Header.Set("Content-Type", gin.MIMEJSON)
			c.Set(gin.BodyBytesKey, []byte(data))
		}
		c.Next()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return len(b), nil
}

func NewEncryptResponseMiddleware(encryptor Encryptor) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()
		// 处理请求
		var response string
		var err error
		response = w.Body.String()
		if len(response) > 0 {
			response, err = encryptor.Encrypt(response)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
		}
		_, err = c.Writer.WriteString(response)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
	}
}

func NewDataEncryptResponseMiddleware(encryptor Encryptor, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()
		// 处理请求
		var (
			obj               JsonResponse
			dataMap           map[string]interface{}
			dataJson, objJson []byte
			response, encStr  string
			err               error
		)
		response = w.Body.String()
		// 加密封装
		err = json.Unmarshal([]byte(response), &obj)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if obj.Data != nil {
			dataJson, err = json.Marshal(obj.Data)
		}
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if len(dataJson) > 0 {
			encStr, err = encryptor.Encrypt(string(dataJson))
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
		}
		if debug {
			if len(dataJson) > 0 {
				err = json.Unmarshal(dataJson, &dataMap)
				if err != nil {
					c.AbortWithStatus(500)
					return
				}
			}
		}
		if dataMap == nil {
			dataMap = make(map[string]interface{})
		}
		dataMap["enc"] = encStr
		obj.Data = dataMap
		objJson, err = json.Marshal(obj)
		_, err = c.Writer.WriteString(string(objJson))
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
	}
}
