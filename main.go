package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"webframe/swf"
)

func main() {
	r := swf.New()
	// 设置静态相对位置
	r.Static("/assets", "./static")
	// 预加载模版
	GetDate := func(t time.Time) string {
		year, month, day := t.Date()
		return fmt.Sprintf("%d-%02d-%02d", year, month, day)
	}
	// 设置模版自定义函数
	r.SetFuncMap(template.FuncMap{
		"GetDate": GetDate,
	})
	// 设置模版数据
	data := map[string]interface{}{
		"Title": "Sherlock WebFramework",
		"User":  "Alex",
		"Items": []string{"Holmes", "Poirot", "Queen", "Dr.Fell"},
		"Time":  time.Date(2025, 8, 13, 0, 0, 0, 0, time.UTC),
	}
	r.LoadHTMLGlob("./static/html/*")
	// 设置routers
	r.GET("/home", func(ctx *swf.Context) {
		ctx.HTML(http.StatusOK, "temp.html", data)
	})
	// GroupTest
	v1 := r.Group("/v1")
	{
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
