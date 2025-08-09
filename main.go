package main

import (
	"net/http"
	"webframe/swf"
)

func main() {
	r := swf.New()
	// 设置routers
	r.GET("/", func(c *swf.Context) {
		c.HTML(http.StatusOK, "<h1>Welcome to Sherlock WebFrameWork!<h1>")
	})
	r.GET("/hello", func(c *swf.Context) {
		c.String(http.StatusOK, "Hi there,%s,you're at %s", c.Query("name"), c.Path)
	})
	r.GET("/login", func(c *swf.Context) {
		c.JSON(http.StatusOK, swf.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	// 启动
	r.Run(":5555")
}
