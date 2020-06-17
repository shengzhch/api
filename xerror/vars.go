package xerror

var (
	Unknown = register(500, 10000, "internal error", "")
	DB      = register(500, 10001, "internal error", "")
	Redis   = register(500, 10002, "internal error", "")
	RPC     = register(500, 10003, "internal error", "")

	NoLogin = register(403, 20000, "用户未登录", "请登录后调用")
)
