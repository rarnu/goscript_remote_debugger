package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqPrint struct {
	DebugToken string `json:"debugToken"`
	Variable   string `json:"variable"`
}

type RespPrint struct {
	Variable string `json:"variable"`
	Output   string `json:"output"`
	Error    string `json:"error"`
}

func ApiPrint(_ *gin.Context, req ReqPrint) (code int, message string, data RespPrint) {
	resp := RespPrint{}
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
		Cmd:     "p",
		StrExpr: req.Variable,
	}
	rp := (<-oc).(RespPrint)
	if rp.Error == "" {
		return 0, "", rp
	} else {
		return global.ERR_VARIABLES, global.ERR_VARIABLES_MESSAGE, rp
	}
}
