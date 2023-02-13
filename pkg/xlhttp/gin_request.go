package xlhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/seanlan/xlvein/pkg/xlerror"
	"net/http"
	"strings"
)

// JsonResponse 返回的json数据格式
type JsonResponse struct {
	Error   int         `json:"error,required"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,required"`
}

type ApiRequest struct {
	ctx *gin.Context
}

func Build(ctx *gin.Context) *ApiRequest {
	return &ApiRequest{ctx}
}

func (r *ApiRequest) RequestParser(args interface{}) (err error) {
	contentType := r.ctx.ContentType()
	method := r.ctx.Request.Method
	switch method {
	case http.MethodPost:
		switch {
		case strings.Contains(contentType, gin.MIMEJSON):
			err = r.ctx.ShouldBindBodyWith(args, binding.JSON)
			break
		case strings.Contains(contentType, gin.MIMEPOSTForm),
			strings.Contains(contentType, gin.MIMEPOSTForm):
			err = r.ctx.MustBindWith(args, binding.Form)
			break
		default:
			err = r.ctx.ShouldBind(args)
		}
		break
	case http.MethodGet:
		err = r.ctx.ShouldBindQuery(args)
	}
	if err != nil {
		r.JsonReturn(xlerror.Wrap(xlerror.ErrRequest, err.Error()))
	}
	return err
}

func (r *ApiRequest) JsonReturn(err error, args ...interface{}) {
	var data interface{}
	if len(args) > 0 {
		data = args[0]
	}
	ec := xlerror.Cause(err)
	if ec.Code() != 0 {
		data = nil
	}
	r.ctx.JSON(http.StatusOK, &JsonResponse{
		Error:   ec.Code(),
		Data:    data,
		Message: ec.Message(),
	})
}

func (r *ApiRequest) GetJWTUID() (int64, error) {
	return r.ctx.GetInt64(JWTIdentityKey), nil
}
