package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqConsoleOutput struct {
	DebugToken string `json:"debugToken"`
}

func ApiConsoleOutput(c *gin.Context, req ReqConsoleOutput) (code int, message string, data []string) {
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
		sarr := global.GlobalSessionPool.GetConsoleOutput(req.DebugToken)
		return 0, "", sarr
	}
	ch := global.GlobalSessionPool.GetCmdChan(req.DebugToken)
	oc := global.GlobalSessionPool.GetOutChan(req.DebugToken)
	ch <- global.DebugCommand{
		Cmd: "out",
	}
	sarr := (<-oc).([]string)
	return 0, "", sarr
}
