//generated by lazy
//author: seanlan

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/seanlan/goutils/xlhttp"
	"github.com/seanlan/xlvein/internal/common/transport"
	"github.com/seanlan/xlvein/internal/config"
	"github.com/seanlan/xlvein/internal/e"
	"github.com/seanlan/xlvein/internal/model"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SocketConnect(c *gin.Context) {
	var (
		err error
	)
	r := xlhttp.Build(c)
	var req model.SocketConnectReq
	err = r.RequestParser(&req)
	if err != nil {
		zap.S().Errorf("request parser error: %s", err.Error())
		return
	}
	app := config.C.GetApp(req.AppID)
	if app == nil {
		zap.S().Errorf("app not found: %s", req.AppID)
		r.JsonReturn(e.ErrAppNotFound)
		return
	}
	jwt := xlhttp.NewJWT(app.AppSecret, 0)
	tag, err := jwt.ParseToken(req.Token)
	if err != nil {
		zap.S().Errorf("ParseToken error: %s", err.Error())
		r.JsonReturn(e.ErrTokenInvalid)
		return
	}
	// 建立websocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("Failed to set websocket upgrade: %#v", err)
		return
	}
	transport.ClientHub.Join(app.AppKey, tag, conn)
	return
}
