package global

import (
	"fmt"
	. "github.com/isyscore/isc-gobase/isc"
	"github.com/rarnu/goscript/dap"
)

type Session struct {
	server *dap.Server
	client *dap.Client
	port   int
	events *ISCList[ISCString]
}

type SessionPool struct {
	pool     ISCMap[string, *Session]
	portList ISCList[int]
}

var GlobalSessionPool = NewSessionPool()

func NewSessionPool() *SessionPool {
	return &SessionPool{
		pool:     ISCMap[string, *Session]{},
		portList: Int(10000, 10100),
	}
}

func (s *SessionPool) Exists(token string) bool {
	_, ok := s.pool[token]
	return ok
}

func (s *SessionPool) GetServer(token string) *dap.Server {
	if ses, ok := s.pool[token]; ok {
		return ses.server
	} else {
		return nil
	}
}

func (s *SessionPool) GetClient(token string) *dap.Client {
	if ses, ok := s.pool[token]; ok {
		return ses.client
	} else {
		return nil
	}
}

func (s *SessionPool) GetEventList(token string) ISCList[ISCString] {
	if ses, ok := s.pool[token]; ok {
		return *ses.events
	} else {
		return nil
	}
}

func (s *SessionPool) AppendEvent(token string, evt string) {
	if ses, ok := s.pool[token]; ok {
		ses.events.Add(ISCString(evt))
	}
}

func (s *SessionPool) Set(token string, svr *dap.Server, cli *dap.Client, port int) {
	if _, ok := s.pool[token]; !ok {
		s.pool[token] = &Session{
			server: svr,
			client: cli,
			port:   port,
			events: &ISCList[ISCString]{},
		}
	}
}

func (s *SessionPool) Remove(token string) {
	if _, ok := s.pool[token]; ok {
		delete(s.pool, token)
	}
}

// GetFreePort 获取一个空闲端口
// 如果没有获取到，表示端口已用完，此时将返回异常
// 目前总共允许开 100 个端口，也就是 100 个 session 同时调试
// PS：为什么使用开端口的形式来进行远程调试，是因为要符合DAP协议，该协议本身只允许单机，
// 并且后续有可能会接入 vscode 或是其他 IDE，因此此处保持每个调试开启新端口的做法
func (s *SessionPool) GetFreePort() (int, error) {
	usedPorts := MapToMapFrom[string, *Session, int](s.pool).Map(func(_ string, ses *Session) int {
		return ses.port
	})
	freePorts := s.portList.Minus(usedPorts)
	if !freePorts.IsEmpty() {
		return freePorts.First(), nil
	} else {
		return -1, fmt.Errorf("调试端口已用完")
	}
}
