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
	server.JsonPost(svr, baseApi+"/init", route.ApiInitDebug)
	server.JsonPost(svr, baseApi+"/exit", route.ApiExitDebug)
	svr.Run(int(port))
}
