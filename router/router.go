package router

import (
	"github.com/gin-gonic/gin"
	"github.com/provider-go/sso/api"
)

type Group struct {
	Router
}

var GroupApp = new(Group)

type Router struct{}

func (s *Router) InitRouter(Router *gin.RouterGroup) {
	{
		// key 注册登录操作
		Router.POST("key", api.LoginByKey)
	}
}
