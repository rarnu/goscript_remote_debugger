package route

import (
	"github.com/gin-gonic/gin"
)

type ReqEvents struct {
	DebugToken string `json:"debugToken"`
}

func ApiEvents(_ *gin.Context, req ReqEvents) (code int, message string, data []map[string]any) {
	/*
		var resp []map[string]any
		if req.DebugToken == "" {
			// 没有传 token，直接退出
			return global.ERR_NO_TOKEN, global.ERR_NO_TOKEN_MESSAGE, resp
		}
		exists := global.GlobalSessionPool.Exists(req.DebugToken)
		if !exists {
			// 没有调试会话，直接退出
			return global.ERR_NO_DEBUG_SESSION, global.ERR_NO_DEBUG_SESSION_MESSAGE, resp
		}
		evtList := global.GlobalSessionPool.GetEventList(req.DebugToken)
		evtJsonList := ListToMapFrom[ISCString, map[string]any](evtList).Map(func(it ISCString) map[string]any {
			var obj map[string]any
			_ = json.Unmarshal([]byte(it), &obj)
			return obj
		})
		return 0, "", evtJsonList

	*/
	return 0, "", nil
}
