package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
	"github.com/google/go-dap"
)

type ReqExitDebug struct {
	DebugToken string `json:"debugToken"`
}

type RespExitDebug struct {
	Exited bool `json:"exited"` // 是否成功退出
}

func ApiExitDebug(c *gin.Context, req ReqExitDebug) (code int, message string, data RespExitDebug) {
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
	// 先退客户端
	cli := global.GlobalSessionPool.GetClient(req.DebugToken)
	disConnReq := &dap.DisconnectRequest{}
	_, err := cli.OnDisconnectRequest(disConnReq)
	if err != nil {
		return global.ERR_EXIT_DEBUG_CLIENT, global.ERR_EXIT_DEBUG_CLIENT_MESSAGE, resp
	}
	// 再退服务端
	svr := global.GlobalSessionPool.GetServer(req.DebugToken)
	svr.Stop()
	// 从会话池中删除记录（此时端口被释放）
	global.GlobalSessionPool.Remove(req.DebugToken)
	resp.Exited = true
	return global.ERR_NONE, global.ERR_NONE_MESSAGE, resp
}
