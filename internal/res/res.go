package res

const Year int64 = 31536000

type FetchType string

const (
	Positive FetchType = "POSITIVE"
	Neutral  FetchType = "NEUTRAL"
	Negative FetchType = "NEGATIVE"
)

type UnitType string

const (
	User UnitType = "USER"
	Chat UnitType = "CHAT"
)

type LocalKeys int

const (
	_ LocalKeys = iota
	CtxAuthData
	CtxUserAgent
	CacheNextRunRegularScheduleAt
)

type LogField string

const (
	RequestID  LogField = "requestID"
	UserID     LogField = "userID"
	SessionKey LogField = "sessionKey"
	Desc       LogField = "desc"
	Loc        LogField = "loc"
)

type LogMsg map[LogField]interface{}
