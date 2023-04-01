package constant

type ctxKey string

const (
	KeyDBCtx      ctxKey = "DB"
	KeyUserIDCtx  ctxKey = "USERID"
	KeyDataSource ctxKey = "DATA_SOURCE"

	SystemID = string("SYSTEM")
	GuestID  = string("GUEST")
)

const (
	SourceDB int = iota
	SourceOS
	SourceRedis
)
