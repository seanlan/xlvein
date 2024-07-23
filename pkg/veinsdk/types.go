package veinsdk

type Response[T any] struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type Empty struct{}

type PushMessageResponse Response[Empty]
