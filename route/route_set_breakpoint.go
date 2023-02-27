package route

import (
	"debugger/global"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ReqSetBreakpoint struct {
	DebugToken string `json:"debugToken"`
	Line       int    `json:"line"`
}

func ApiSetBreakpoint(_ *gin.Context, req ReqSetBreakpoint) (code int, message string, data CommonResp) {
	resp := CommonResp{}
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
		Cmd:     "sb",
		IntExpr: req.Line,
	}
	err := <-oc
	if err == nil {
		return 0, "", resp
	} else {
		return global.ERR_SET_BREAKPOINT, fmt.Sprintf(global.ERR_SET_BREAKPOINT_MESSAGE, err), resp
	}
}
