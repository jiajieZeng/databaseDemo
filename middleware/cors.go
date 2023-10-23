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
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

// func Cors() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         method := c.Request.Method
//         origin := c.Request.Header.Get("Origin") //请求头部
//         if origin != "" {
//             //接收客户端发送的origin （重要！）
//             c.Writer.Header().Set("Access-Control-Allow-Origin", origin) 
//             //服务器支持的所有跨域请求的方法
//             c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") 
//             //允许跨域设置可以返回其他子段，可以自定义字段
//             c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
//             // 允许浏览器（客户端）可以解析的头部 （重要）
//             c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers") 
//             //设置缓存时间
//             c.Header("Access-Control-Max-Age", "172800") 
//             //允许客户端传递校验信息比如 cookie (重要)
//             c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                                                                                          
//         }

//         //允许类型校验 
//         if method == "OPTIONS" {
//             c.JSON(http.StatusOK, "ok!")
//         }
//         c.Next()
//     }
// }