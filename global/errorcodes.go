package global

const (
	ERR_NONE, ERR_NONE_MESSAGE                           = 0, ""
	ERR_NO_TOKEN, ERR_NO_TOKEN_MESSAGE                   = 1, "没有提供会话 Token"
	ERR_NO_FREE_PORT, ERR_NO_FREE_PORT_MESSAGE           = 2, "没有空闲的调试端口"
	ERR_INIT_FAILED, ERR_INIT_FAILED_MESSAGE             = 3, "初始化远程调试器失败"
	ERR_NO_DEBUG_SESSION, ERR_NO_DEBUG_SESSION_MESSAGE   = 4, "没有相应的调试会话"
	ERR_EXIT_DEBUG_CLIENT, ERR_EXIT_DEBUG_CLIENT_MESSAGE = 5, "退出调试客户端失败"
	ERR_LAUNCH_DEBUG, ERR_LAUNCH_DEBUG_MESSAGE           = 6, "启动调试失败 %v"
	ERR_SET_BREAKPOINT, ERR_SET_BREAKPOINT_MESSAGE       = 7, "设置断点失败 %v"
	ERR_CLEAR_BREAKPOINT, ERR_CLEAR_BREAKPOINT_MESSAGE   = 7, "清除断点失败 %v"
	ERR_GET_BREAKPOINTS, ERR_GET_BREAKPOINTS_MESSAGE     = 7, "获取断点列表失败"
	ERR_CONTINUE, ERR_CONTINUE_MESSAGE                   = 8, "继续执行到下一个断点失败 %v"
	ERR_NEXT, ERR_NEXT_MESSAGE                           = 9, "执行到下一行失败 %v"
	ERR_VARIABLES, ERR_VARIABLES_MESSAGE                 = 10, "获取变量值失败"
	ERR_EVALUATE, ERR_EVALUATE_MESSAGE                   = 11, "表达式求值失败 %v"
	ERR_DEBUG_IS_FINISHED, ERR_DEBUG_IS_FINISHED_MESSAGE = 12, "该调试已结束"
	ERR_LAST_LINE, ERR_LAST_LINE_MESSAGE                 = 0, "已是最后一行，调试完毕"
	ERR_NO_EXPORT, ERR_NO_EXPORT_MESSAGE                 = 13, "调试还未结束，不能获取返回值"
)
