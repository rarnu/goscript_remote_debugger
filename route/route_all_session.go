package route

import (
	"debugger/global"
	"github.com/gin-gonic/gin"
)

type ReqAllSession struct {
}

func ApiAllSession(_ *gin.Context, _ ReqAllSession) (code int, message string, data map[string]any) {
	return 0, "", global.GlobalSessionPool.AllSession()
}
