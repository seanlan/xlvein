package e

import "github.com/seanlan/goutils/xlerror"

var (
	ErrTokenInvalid = xlerror.ErrToken
	ErrAppNotFound  = xlerror.NewError(1001, "app not found")
)
