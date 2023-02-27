package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqExport struct {
	DebugToken string `json:"debugToken"`
}

func ApiExport(_ *gin.Context, req ReqExport) (code int, message string, data any) {
	if req.DebugToken == "" {
		// 没有传 token，直接退出
		return global.ERR_NO_TOKEN, global.ERR_NO_TOKEN_MESSAGE, nil
	}
	exists := global.GlobalSessionPool.Exists(req.DebugToken)
	if !exists {
		// 没有调试会话，直接退出
		return global.ERR_NO_DEBUG_SESSION, global.ERR_NO_DEBUG_SESSION_MESSAGE, nil
	}
	if global.GlobalSessionPool.IsClosed(req.DebugToken) {
		// 如果已关闭，从池子拿
		ex := global.GlobalSessionPool.GetExport(req.DebugToken)
		return 0, "", ex
	}
	return global.ERR_NO_EXPORT, global.ERR_NO_EXPORT_MESSAGE, nil
}
