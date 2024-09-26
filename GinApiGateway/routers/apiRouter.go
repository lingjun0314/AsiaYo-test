package routers

import (
	"ginapigateway/controllers"

	"github.com/gin-gonic/gin"
)

var validation *controllers.ValidatoinController

func InitApi(r *gin.Engine) {
	r.POST("/api/orders", validation.CheckDataType)
}
