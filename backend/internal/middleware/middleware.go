package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// CORS 跨域中间件
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Auth(r *ghttp.Request) {
	walletAddress := r.GetHeader("X-Wallet-Address")
	if walletAddress == "" {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "Unauthorized: wallet address required",
		})
		return
	}
	r.SetCtxVar("walletAddress", walletAddress)
	r.Middleware.Next()
}

// RequestLog 请求日志中间件
func RequestLog(r *ghttp.Request) {
	g.Log().Info(r.Context(), "Request:",
		"method", r.Method,
		"path", r.URL.Path,
		"clientIP", r.GetClientIp(),
		"userAgent", r.UserAgent(),
	)
	r.Middleware.Next()
}

// ErrorHandler 错误处理中间件
func ErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		r.Response.ClearBuffer()
		r.Response.WriteJson(g.Map{
			"code":    500,
			"message": err.Error(),
		})
	}
}
