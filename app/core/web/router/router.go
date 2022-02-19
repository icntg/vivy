package router

import (
	"app/core/web/api/utility"
	"github.com/gin-gonic/gin"
)

func WebRouters(engine *gin.Engine) {
	api := engine.Group("api")
	{
		g := api.Group("utility")
		{
			g.GET("/captcha", utility.NewCaptcha)
		}
	}

}
