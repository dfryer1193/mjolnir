package middleware

type ctxKey int

const (
	errorCtxKey ctxKey = iota
	requestIDKey
)
