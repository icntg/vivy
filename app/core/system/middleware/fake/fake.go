package fake

import (
	"github.com/gin-gonic/gin"
)

func GinFake() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Server", "Apache/2.2.15 (CentOS)")
		ctx.Header("X-Powered-By", "PHP/5.2.17")
		ctx.Next()
	}
}
