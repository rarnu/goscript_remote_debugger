package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqEvaluate struct {
	DebugToken string `json:"debugToken"`
	Expression string `json:"expression"`
}

type RespEvaluate struct {
	Expression string `json:"expression"`
	Error      string `json:"error"`
	Value      any    `json:"value"`
}

func ApiEvaluate(_ *gin.Context, req ReqEvaluate) (code int, message string, data RespEvaluate) {
	resp := RespEvaluate{}
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
		Cmd:     "e",
		StrExpr: req.Expression,
	}
	ex := (<-oc).(RespEvaluate)
	if ex.Error != "" {
		ex.Error += ", 调试器已退出"
	}
	return 0, "", ex

	/*
		if req.DebugToken == "" {
			// 没有传 token，直接退出
			return global.ERR_NO_TOKEN, global.ERR_NO_TOKEN_MESSAGE, resp
		}
		exists := global.GlobalSessionPool.Exists(req.DebugToken)
		if !exists {
			// 没有调试会话，直接退出
			return global.ERR_NO_DEBUG_SESSION, global.ERR_NO_DEBUG_SESSION_MESSAGE, resp
		}
		cli := global.GlobalSessionPool.GetClient(req.DebugToken)
		evaReq := &godap.EvaluateRequest{
			Arguments: godap.EvaluateArguments{
				Expression: req.Expression,
			},
		}
		evaResp, err := cli.Evaluate(evaReq)
		if err != nil {
			return global.ERR_EVALUATE, fmt.Sprintf(global.ERR_EVALUATE_MESSAGE, err), resp
		}
		return 0, "", *evaResp

	*/
}
