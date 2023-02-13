package e

import "github.com/seanlan/goutils/xlerror"

var (
	ErrTokenInvalid   = xlerror.ErrToken
	ErrAppNotFound    = xlerror.NewError(1001, "app not found")
	ErrSignInvalid    = xlerror.NewError(1002, "sign invalid")
	ErrMessageInvalid = xlerror.NewError(1003, "message invalid")
)
