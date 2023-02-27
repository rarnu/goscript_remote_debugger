package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqGetBreakpoints struct {
	DebugToken string `json:"debugToken"`
}

func ApiGetBreakpoints(_ *gin.Context, req ReqGetBreakpoints) (code int, message string, data []int) {
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
		return global.ERR_DEBUG_IS_FINISHED, global.ERR_DEBUG_IS_FINISHED_MESSAGE, nil
	}

	ch := global.GlobalSessionPool.GetCmdChan(req.DebugToken)
	oc := global.GlobalSessionPool.GetOutChan(req.DebugToken)
	ch <- global.DebugCommand{
		Cmd: "b",
	}
	mp := <-oc
	if mp != nil {
		return 0, "", mp.(map[string][]int)["index.js"]
	} else {
		return global.ERR_GET_BREAKPOINTS, global.ERR_GET_BREAKPOINTS_MESSAGE, nil
	}
}
