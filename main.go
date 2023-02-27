package main

import (
	"debugger/route"
	"debugger/server"
	"os"
	"strconv"
)

func printHelp() {

}

func main() {
	if len(os.Args) != 2 {
		printHelp()
		return
	}

	port, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		printHelp()
		return
	}

	baseApi := "/api/goscript/debug"

	svr := server.New()

	// 获取所有调试器(输出 code, running)
	server.JsonGet(svr, baseApi+"/allSession", route.ApiAllSession)
	// 开始调试（传入脚本内容，参数内容）
	server.JsonPost(svr, baseApi+"/launch", route.ApiLaunchDebug)
	// 设置断点
	server.JsonPost(svr, baseApi+"/setBreakpoint", route.ApiSetBreakpoint)
	// 清除断点
	server.JsonPost(svr, baseApi+"/clearBreakpoint", route.ApiClearBreakpoint)
	// 获取断点列表
	server.JsonPost(svr, baseApi+"/getBreakpoints", route.ApiGetBreakpoints)
	// 执行到下一个断点（或执行到结束处）	// continue
	server.JsonPost(svr, baseApi+"/continue", route.ApiContinue)
	// 执行到下一行代码（或执行到结束处）	// next
	server.JsonPost(svr, baseApi+"/next", route.ApiNext)
	// 获取当前代码的上下文
	server.JsonPost(svr, baseApi+"/list", route.ApiList)
	// 获取调试堆栈
	server.JsonPost(svr, baseApi+"/backtrace", route.ApiBacktrace)
	// 获取变量值 // variables
	server.JsonPost(svr, baseApi+"/print", route.ApiPrint)
	// 计算表达式	// evaluate
	server.JsonPost(svr, baseApi+"/evaluate", route.ApiEvaluate)
	// 获取控制台输出
	server.JsonPost(svr, baseApi+"/consoleOutput", route.ApiConsoleOutput)
	// 获取脚本的整体返回值
	server.JsonPost(svr, baseApi+"/export", route.ApiExport)
	// 退出调试（关闭调试器）
	server.JsonPost(svr, baseApi+"/exit", route.ApiExitDebug)
	svr.Run(int(port))
}
