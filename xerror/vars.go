package xerror

var (
	Unknown = register(10000, "internal error", "")
	DB      = register(10001, "internal error", "")
	Redis   = register(10002, "internal error", "")
	RPC     = register(10003, "internal error", "")

	NoLogin = register(20000, "用户未登录", "请登录后调用")
)
