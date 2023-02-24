package server

import (
	"debugger/server/plugin"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FuncGinRoute[REQUEST any, RESPONSE any] func(c *gin.Context, req REQUEST) (code int, message string, data RESPONSE)

type Server struct {
	engine  *gin.Engine
	methods map[string]func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

func New() *Server {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	eng.Use(gin.Recovery(), plugin.Cors())
	return &Server{
		engine: eng,
		methods: map[string]func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes{
			"GET":     eng.GET,
			"POST":    eng.POST,
			"PUT":     eng.PUT,
			"DELETE":  eng.DELETE,
			"OPTIONS": eng.OPTIONS,
			"HEAD":    eng.HEAD,
			"PATCH":   eng.PATCH,
		},
	}
}

func (s *Server) Use(plugins ...gin.HandlerFunc) {
	s.engine.Use(plugins...)
}

func (s *Server) Run(port int) {
	_ = s.engine.Run(fmt.Sprintf(":%d", port))
}

func JsonGet[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "GET", path, route)
}

func JsonPost[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "POST", path, route)
}

func JsonPut[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "PUT", path, route)
}

func JsonDelete[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "DELETE", path, route)
}

func JsonOptions[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "OPTIONS", path, route)
}

func JsonHead[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "HEAD", path, route)
}

func JsonPatch[REQUEST any, RESPONSE any](s *Server, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	jsonHandler(s, "PATCH", path, route)
}

func jsonHandler[REQUEST any, RESPONSE any](s *Server, method string, path string, route FuncGinRoute[REQUEST, RESPONSE]) {
	s.methods[method](path, func(ctx *gin.Context) {
		var _req REQUEST
		switch method {
		case "GET", "DELETE":
			_query := map[string][]string{}
			_ = ctx.ShouldBindQuery(&_query)
			b, _ := json.Marshal(mapStringArrToMapAny(_query))
			_ = json.Unmarshal(b, &_req)
		case "POST", "PUT", "PATCH":
			_ = ctx.ShouldBindJSON(&_req)
		}
		_code, _message, _data := route(ctx, _req)
		ctx.JSON(http.StatusOK, &BaseResponse[any]{Code: _code, Message: _message, Data: _data})
	})
}

func mapStringArrToMapAny(m map[string][]string) map[string]any {
	ret := map[string]any{}
	for k, v := range m {
		if v == nil || len(v) == 0 {
			ret[k] = ""
		} else {
			ret[k] = v[0]
		}
	}
	return ret
}
