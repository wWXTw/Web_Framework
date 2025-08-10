package main

import (
	"net/http"
	"webframe/swf"
)

func main() {
	r := swf.New()
	// 设置routers
	r.GET("/", func(ctx *swf.Context) {
		ctx.HTML(http.StatusOK, "<h1>Welcome to Sherlock WebFrameWork!<h1>")
	})
	r.GET("/hello", func(ctx *swf.Context) {
		ctx.String(http.StatusOK, "Hi there,%s,you're at %s\n", ctx.Query("name"), ctx.Path)
	})
	r.GET("/hello/:name", func(ctx *swf.Context) {
		ctx.String(http.StatusOK, "Hi there,%s,you're at %s\n", ctx.Param("name"), ctx.Path)
	})
	r.GET("/assets/*filepath", func(ctx *swf.Context) {
		ctx.JSON(http.StatusOK, swf.H{
			"filepath": ctx.Param("filepath"),
		})
	})
	r.GET("/login", func(ctx *swf.Context) {
		ctx.JSON(http.StatusOK, swf.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})
	// 启动
	r.Run(":5555")
}
