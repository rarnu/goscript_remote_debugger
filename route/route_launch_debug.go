package route

import (
	"bytes"
	"debugger/global"
	"debugger/util"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/rarnu/goscript"
	"github.com/rarnu/goscript/module/console"
	"github.com/rarnu/goscript/module/require"
)

type ReqLaunchDebug struct {
	// DebugToken string         `json:"debugToken"`
	Code   string         `json:"code"`
	Params map[string]any `json:"params"`
}

type RespLaunchDebug struct {
	DebugToken string `json:"debugToken"`
	Reason     string `json:"reason"`
}

func ApiLaunchDebug(_ *gin.Context, req ReqLaunchDebug) (code int, message string, data RespLaunchDebug) {
	resp := RespLaunchDebug{}

	token := goid.GenerateUUID()
	dbgCode := "debugger\n" + req.Code

	filename := "index.js"
	content := util.GenerateSourceMap(filename, dbgCode)

	printer := &console.ExecPrinter{}
	runtime := goscript.New()
	dbg := runtime.AttachDebugger()
	registry := new(require.Registry)
	registry.Enable(runtime)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))

	for k, v := range req.Params {
		_ = runtime.Set(k, v)
	}

	ch := make(chan global.DebugCommand)
	global.GlobalSessionPool.Set(token, dbgCode, runtime, dbg, ch)

	chReason := make(chan string)
	go func(rtm *goscript.Runtime, c chan global.DebugCommand, cr chan string, d *goscript.Debugger, prt *console.ExecPrinter, fn string, tkn string) {
		defer d.Detach()
		reason := d.Continue()
		cr <- string(reason)
		for {
			cmd := <-c
			switch cmd.Cmd {
			case "exit":
				break
			case "c": // contiue
				reason = d.Continue()
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- reason
			case "n": // next
				err := d.Next()
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- err
			case "sb": // setBreakpoint
				err := d.SetBreakpoint(fn, cmd.IntExpr)
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- err
			case "b": // breakpoints
				mp, _ := d.Breakpoints()
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- mp
			case "cb": // clearBreakpoint
				err := d.ClearBreakpoint(fn, cmd.IntExpr)
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- err
			case "l": // list
				lines, _ := d.List()
				currentLine := dbg.Line()
				lineIndex := currentLine - 1
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- RespList{
					Lines:     lines,
					LineIndex: lineIndex,
				}
			case "bt": // backtrace
				stack := rtm.CaptureCallStack(0, nil)
				var backtrace bytes.Buffer
				backtrace.WriteRune('\n')
				for _, frame := range stack {
					frame.Write(&backtrace)
					backtrace.WriteRune('\n')
				}
				str := backtrace.String()
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- str
			case "out": // output
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- prt.Lines
			case "p": // print
				str, err := d.Print(cmd.StrExpr)
				errStr := ""
				if err != nil {
					errStr = err.Error()
				}
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- RespPrint{
					Variable: cmd.StrExpr,
					Output:   str,
					Error:    errStr,
				}
			case "e": // evaluate
				v, err := d.Exec(cmd.StrExpr)
				errStr := ""
				if err != nil {
					errStr = err.Error()
				}
				var vRet any
				if v != nil {
					vRet = v.Export()
				}
				oc := global.GlobalSessionPool.GetOutChan(tkn)
				oc <- RespEvaluate{
					Expression: cmd.StrExpr,
					Error:      errStr,
					Value:      vRet,
				}
				if err != nil {
					// 如果 evaluate 有异常，直接把调试器退掉，因为此时状态已经被变更，无法恢复
					d.Detach()
					global.GlobalSessionPool.SetStatus(tkn, nil, err, false, prt.Lines)
					break
				}
			}
		}
	}(runtime, ch, chReason, dbg, printer, filename, token)

	go func(r *goscript.Runtime, fn string, cnt []byte, tkn string, ch chan global.DebugCommand, prt *console.ExecPrinter) {
		console.Enable(r)
		exp, err := r.RunScript(fn, string(cnt))
		close(ch)
		if !global.GlobalSessionPool.IsClosed(tkn) {
			oc := global.GlobalSessionPool.GetOutChan(tkn)
			oc <- "finished"
			global.GlobalSessionPool.SetStatus(tkn, exp.Export(), err, false, prt.Lines)
		}
	}(runtime, filename, content, token, ch, printer)
	resp.DebugToken = token
	resp.Reason = <-chReason
	close(chReason)
	return 0, "", resp
}
