package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqList struct {
	DebugToken string `json:"debugToken"`
}

type RespList struct {
	Lines     []string `json:"lines"`
	LineIndex int      `json:"lineIndex"`
}

func ApiList(_ *gin.Context, req ReqList) (code int, message string, data RespList) {
	resp := RespList{}
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
		Cmd: "l",
	}
	lo := (<-oc).(RespList)
	return 0, "", lo
}
