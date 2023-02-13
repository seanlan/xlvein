//generated by lazy
//author: seanlan

package service

import (
	"context"
	"encoding/json"
	"github.com/seanlan/xlvein/app/common/exchange"
	"github.com/seanlan/xlvein/app/common/transport"
	"github.com/seanlan/xlvein/app/dao"
	"github.com/seanlan/xlvein/app/dao/sqlmodel"
	"github.com/seanlan/xlvein/app/e"
	"github.com/seanlan/xlvein/app/model"
	"github.com/seanlan/xlvein/pkg/veinsdk"
	"time"
)

func PushMessage(ctx context.Context, req model.PushMessageReq) (resp model.PushMessageResp, err error) {
	var (
		appQ = sqlmodel.ApplicationColumns
		app  sqlmodel.Application
	)
	err = dao.FetchApplication(ctx, &app, appQ.AppKey.Eq(req.AppKey))
	if err != nil || app.ID == 0 {
		err = e.ErrAppNotFound
		return
	}
	sdk := veinsdk.New("", req.AppKey, app.AppSecret)
	if sdk.GetSign(map[string]interface{}{
		"app_key": app.AppKey,
		"send_to": req.SendTo,
		"message": req.Message,
		"nonce":   req.Nonce,
	}) != req.Sign {
		err = e.ErrSignInvalid
		return
	}
	var msg map[string]interface{}
	err = json.Unmarshal([]byte(req.Message), &msg)
	if err != nil {
		err = e.ErrMessageInvalid
		return
	}
	transport.ClientHub.PushToExchange(
		req.AppKey,
		transport.Message{
			From:   "",
			To:     req.SendTo,
			Event:  exchange.EventSystem,
			Data:   msg,
			SendAt: time.Now().UnixMilli(),
		})
	return
}
