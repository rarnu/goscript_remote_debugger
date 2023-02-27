package global

import (
	. "github.com/isyscore/isc-gobase/isc"
	"github.com/rarnu/goscript"
)

type DebugCommand struct {
	Cmd     string
	IntExpr int
	StrExpr string
}

type Session struct {
	// server *dap.Server
	// client *dap.Client
	// port   int
	// events *ISCList[ISCString]

	Code     string
	Runtime  *goscript.Runtime
	Debugger *goscript.Debugger

	ConsoleLines []string
	Export       any
	Running      bool
	CmdChan      chan DebugCommand
	OutChan      chan any
	Error        error
}

type SessionPool struct {
	pool ISCMap[string, *Session]
	// portList ISCList[int]
}

var GlobalSessionPool = NewSessionPool()

func NewSessionPool() *SessionPool {
	return &SessionPool{
		pool: ISCMap[string, *Session]{},
		// portList: Int(10000, 10100),
	}
}

func (s *SessionPool) Exists(token string) bool {
	_, ok := s.pool[token]
	return ok
}

func (s *SessionPool) GetCmdChan(token string) chan DebugCommand {
	if ses, ok := s.pool[token]; ok {
		return ses.CmdChan
	} else {
		return nil
	}
}

func (s *SessionPool) GetOutChan(token string) chan any {
	if ses, ok := s.pool[token]; ok {
		return ses.OutChan
	} else {
		return nil
	}
}

func (s *SessionPool) IsClosed(token string) bool {
	if ses, ok := s.pool[token]; ok {
		return !ses.Running
	} else {
		return false
	}
}

func (s *SessionPool) GetConsoleOutput(token string) []string {
	if ses, ok := s.pool[token]; ok {
		return ses.ConsoleLines
	} else {
		return nil
	}
}

func (s *SessionPool) GetExport(token string) any {
	if ses, ok := s.pool[token]; ok {
		return ses.Export
	} else {
		return nil
	}
}

//func (s *SessionPool) GetServer(token string) *dap.Server {
//	if ses, ok := s.pool[token]; ok {
//		return ses.server
//	} else {
//		return nil
//	}
//}

//func (s *SessionPool) GetClient(token string) *dap.Client {
//	if ses, ok := s.pool[token]; ok {
//		return ses.client
//	} else {
//		return nil
//	}
//}

//func (s *SessionPool) GetEventList(token string) ISCList[ISCString] {
//	if ses, ok := s.pool[token]; ok {
//		return *ses.events
//	} else {
//		return nil
//	}
//}

//func (s *SessionPool) AppendEvent(token string, evt string) {
//	if ses, ok := s.pool[token]; ok {
//		ses.events.Add(ISCString(evt))
//	}
//}

//func (s *SessionPool) Set(token string, svr *dap.Server, cli *dap.Client, port int) {
//	if _, ok := s.pool[token]; !ok {
//		s.pool[token] = &Session{
//			server: svr,
//			client: cli,
//			port:   port,
//			events: &ISCList[ISCString]{},
//		}
//	}
//}

func (s *SessionPool) Set(token string, code string, r *goscript.Runtime, d *goscript.Debugger, ch chan DebugCommand) {
	if _, ok := s.pool[token]; !ok {
		s.pool[token] = &Session{
			Code:     code,
			Runtime:  r,
			Debugger: d,
			Running:  true,
			CmdChan:  ch,
			OutChan:  make(chan any),
		}
	}
}

func (s *SessionPool) SetStatus(token string, exp any, err error, running bool, consoleLines []string) {
	if ses, ok := s.pool[token]; ok {
		ses.Export = exp
		ses.Error = err
		ses.Running = running
		ses.ConsoleLines = consoleLines
	}
}

func (s *SessionPool) Remove(token string) {
	if ses, ok := s.pool[token]; ok {
		close(ses.OutChan)
		delete(s.pool, token)
	}
}

// GetFreePort 获取一个空闲端口
// 如果没有获取到，表示端口已用完，此时将返回异常
// 目前总共允许开 100 个端口，也就是 100 个 session 同时调试
// PS：为什么使用开端口的形式来进行远程调试，是因为要符合DAP协议，该协议本身只允许单机，
// 并且后续有可能会接入 vscode 或是其他 IDE，因此此处保持每个调试开启新端口的做法
//func (s *SessionPool) GetFreePort() (int, error) {
//	usedPorts := MapToMapFrom[string, *Session, int](s.pool).Map(func(_ string, ses *Session) int {
//		return ses.port
//	})
//	freePorts := s.portList.Minus(usedPorts)
//	if !freePorts.IsEmpty() {
//		return freePorts.First(), nil
//	} else {
//		return -1, fmt.Errorf("调试端口已用完")
//	}
//}

func (s *SessionPool) AllSession() map[string]any {
	m := map[string]any{}
	for k, v := range s.pool {
		m[k] = map[string]any{
			"code":      v.Code,
			"isRunning": v.Running,
		}
	}
	return m
}
