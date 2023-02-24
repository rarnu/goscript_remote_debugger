package plugin

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
)

func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method
        origin := c.Request.Header.Get("Origin")
        if origin != "" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
            c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
            c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
            c.Header("Access-Control-Max-Age", "172800")
            c.Header("Access-Control-Allow-Credentials", "true")
        }
        //允许类型校验
        if method == "OPTIONS" {
            c.JSON(http.StatusOK, "ok!")
        }
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("跨域中断: %v\n", err)
            }
        }()
        c.Next()
    }
}
