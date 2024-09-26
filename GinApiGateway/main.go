package main

import (
	"ginapigateway/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routers.InitApi(r)
	r.Run()
}