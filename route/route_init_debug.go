package route

import (
	"debugger/global"
	"fmt"
	"github.com/gin-gonic/gin"
	godap "github.com/google/go-dap"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/rarnu/goscript/dap"
)

// 开始一个新的调试

type ReqInitDebug struct {
}

type RespInitDebug struct {
	DebugToken   string `json:"debugToken"`   // 调试token，后续的请求都必须传入，以供识别调试的session
	Capabilities any    `json:"capabilities"` // 初始化装载信息
}

func ApiInitDebug(c *gin.Context, req ReqInitDebug) (code int, message string, data RespInitDebug) {
	resp := RespInitDebug{}
	dispatchPort, err := global.GlobalSessionPool.GetFreePort()
	if err != nil {
		return global.ERR_NO_FREE_PORT, global.ERR_NO_FREE_PORT_MESSAGE, resp
	}
	// 生成调试 token
	token := goid.GenerateUUID()

	// 启动调试服务
	svr := dap.StartInstance(dispatchPort)
	// 初始化调试客户端
	cli, _ := dap.NewClient(fmt.Sprintf("127.0.0.1:%d", dispatchPort))
	cli.MsgChan = make(chan string)
	global.GlobalSessionPool.Set(token, svr, cli, dispatchPort)
	// 此处接收事件，作为 DAP 协议的 EventBus 而存在
	go func(c *dap.Client, tk string) {
		for {
			msg := <-c.MsgChan
			global.GlobalSessionPool.AppendEvent(tk, msg)
		}
	}(cli, token)

	// 请求初始化
	initReq := &godap.InitializeRequest{
		Arguments: godap.InitializeRequestArguments{
			PathFormat:      "path",
			LinesStartAt1:   true,
			ColumnsStartAt1: true,
		},
	}
	initResp, err := cli.Initialize(initReq)
	if err != nil {
		// 初始化出错的情况，将原先写入的 session 信息删除
		global.GlobalSessionPool.Remove(token)
		return global.ERR_INIT_FAILED, global.ERR_INIT_FAILED_MESSAGE, resp
	}
	resp.DebugToken = token
	resp.Capabilities = initResp.Body
	return 0, "", resp
}
