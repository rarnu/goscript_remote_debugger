package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqExitDebug struct {
	DebugToken string `json:"debugToken"`
}

type RespExitDebug struct {
	Exited bool `json:"exited"` // 是否成功退出
}

func ApiExitDebug(_ *gin.Context, req ReqExitDebug) (code int, message string, data RespExitDebug) {
	resp := RespExitDebug{Exited: false}
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
		// 如果已经关闭了（是调试到终点后关闭）
		if global.GlobalSessionPool.Exists(req.DebugToken) {
			// 从池子里移除调试会话
			global.GlobalSessionPool.Remove(req.DebugToken)
			resp.Exited = true
			return 0, "", resp
		}
		// 如果已经移除过了，抛个异常
		return global.ERR_DEBUG_IS_FINISHED, global.ERR_DEBUG_IS_FINISHED_MESSAGE, resp
	}
	// 如果还没有关闭（在调试中需要强制关闭）
	ch := global.GlobalSessionPool.GetCmdChan(req.DebugToken)
	ch <- global.DebugCommand{
		Cmd: "exit",
	}
	global.GlobalSessionPool.Remove(req.DebugToken)
	resp.Exited = true
	return 0, "", resp
}
