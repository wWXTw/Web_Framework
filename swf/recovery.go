package swf

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 打印错误数据
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery中间件函数
func Recovery() HandleFunc {
	return func(ctx *Context) {
		defer func() {
			// 发生panic后进行recover
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n", trace(message))
				ctx.Status(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}
