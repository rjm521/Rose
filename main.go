package main

import (
	"RunCodeServer/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", server.Ping)
	r.GET("/runcode", server.HandleCode)
	r.Run(":8090")
}