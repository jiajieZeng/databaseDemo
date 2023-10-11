package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*

关于跨域的解决方法，大部分可以分为 2 种

nginx反向代理解决跨域
服务端设置Response Header(响应头部)的Access-Control-Allow-Origin
对于后端开发来说，第 2 种的操作性更新灵活，这里也讲一下 Gin 是如何做到的

在 Gin 中提供了 middleware (中间件) 来做到在一个请求前后处理响应的逻辑，
这里我们使用中间来做到在每次请求是添加上 Access-Control-Allow-Origin 头部

*/

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
