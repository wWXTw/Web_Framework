package main

import (
	"net/http"
	"webframe/swf"
)

func main() {
	r := swf.New()
	// 设置routers
	r.GET("/home", func(ctx *swf.Context) {
		ctx.HTML(http.StatusOK, "<h1>Welcome to Sherlock WebFrameWork!<h1>")
	})
	// GroupTest
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(ctx *swf.Context) {
			ctx.HTML(http.StatusOK, "<h1>Sherlock WebFramework V1<h1>")
		})
		v1.GET("/hello", func(ctx *swf.Context) {
			ctx.String(http.StatusOK, "Hi there,%s,you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}
	v2 := r.Group("/v2")
	// v2组使用Logger中间件
	v2.Use(swf.Logger())
	{
		v2.GET("/hello/:name", func(ctx *swf.Context) {
			ctx.String(http.StatusOK, "Hi there,%s,you're at %s\n", ctx.Param("name"), ctx.Path)
		})
		v2.GET("/assets/*filepath", func(ctx *swf.Context) {
			ctx.JSON(http.StatusOK, swf.H{
				"filepath": ctx.Param("filepath"),
			})
		})
		v2.GET("/login", func(ctx *swf.Context) {
			ctx.JSON(http.StatusOK, swf.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
			})
		})
	}
	// 启动
	r.Run(":5555")
}
