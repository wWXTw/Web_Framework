package swf

import (
	"log"
	"time"
)

// 一个用来记录并打印时间信息的中间件函数
func Logger() HandleFunc {
	return func(ctx *Context) {
		// 记录开始时间
		t := time.Now()
		// 执行Next
		ctx.Next()
		// 打印信息
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
