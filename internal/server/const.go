package server

const (
	SessionKeyUser = "__USER__"
	TOKEN          = "token"
	GUEST          = "__guest__"
	AUTH_CONTEXT   = "__auth_context__"

	Relam             = "Wikinote"
	AuthenicateString = `Basic realm="` + Relam + `"`
)
