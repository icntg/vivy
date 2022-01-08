package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUid(c *gin.Context) interface{} {
	session := sessions.Default(c)
	//v := session.Get(conf.Cfg.Token)
	// TODO:
	v := session.Get("what's_this?")
	if v == nil {
		return nil
	}
	return session.Get(v)
}
