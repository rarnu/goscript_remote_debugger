package route

import (
	"debugger/global"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ReqContinue struct {
	DebugToken string `json:"debugToken"`
}

type RespContinue struct {
	Reason string `json:"reason"`
}

func ApiContinue(_ *gin.Context, req ReqContinue) (code int, message string, data RespContinue) {
	resp := RespContinue{}
	if req.DebugToken == "" {
		// 没有传 token，直接退出
		return global.ERR_NO_TOKEN, global.ERR_NO_TOKEN_MESSAGE, resp
	}
	exists := global.GlobalSessionPool.Exists(req.DebugToken)
	if !exists {
		// 没有调试会话，直接退出
		return global.ERR_NO_DEBUG_SESSION, global.ERR_NO_DEBUG_SESSION_MESSAGE, resp
	}
	if global.GlobalSessionPool.IsClosed(req.DebugToken) {
		return global.ERR_DEBUG_IS_FINISHED, global.ERR_DEBUG_IS_FINISHED_MESSAGE, resp
	}
	ch := global.GlobalSessionPool.GetCmdChan(req.DebugToken)
	oc := global.GlobalSessionPool.GetOutChan(req.DebugToken)
	ch <- global.DebugCommand{
		Cmd: "c",
	}
	rs := <-oc
	if s, ok := rs.(string); ok {
		if s == "finished" {
			return global.ERR_LAST_LINE, global.ERR_LAST_LINE_MESSAGE, resp
		}
	}
	resp.Reason = fmt.Sprintf("%v", rs)
	return 0, "", resp
}
