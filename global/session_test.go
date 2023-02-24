package global

import "testing"

func TestAppendEvent(t *testing.T) {
	token := "abcdefg"
	GlobalSessionPool.Set(token, nil, nil, 0)
	GlobalSessionPool.AppendEvent(token, "test event 1")
	GlobalSessionPool.AppendEvent(token, "test event 2")
	list := GlobalSessionPool.GetEventList(token)
	t.Logf("list = %+v\n", list)
}

func TestFreePort(t *testing.T) {
	port1, _ := GlobalSessionPool.GetFreePort()
	t.Logf("port1 = %d\n", port1) // 10000
	GlobalSessionPool.Set("a", nil, nil, 10000)
	GlobalSessionPool.Set("b", nil, nil, 10001)
	GlobalSessionPool.Set("c", nil, nil, 10003)
	port2, _ := GlobalSessionPool.GetFreePort()
	t.Logf("port2 = %d\n", port2) // 10002
	GlobalSessionPool.Set("d", nil, nil, 10002)
	port3, _ := GlobalSessionPool.GetFreePort()
	t.Logf("port3 = %d\n", port3) // 10004
}
