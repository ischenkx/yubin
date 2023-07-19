package rest

import (
	"github.com/gin-gonic/gin"
	"log"
	"yubin/lib/api/rest/pubsub"
	"yubin/lib/api/rest/sources"
	"yubin/lib/api/rest/templates"
	"yubin/lib/api/rest/users"
	yubin "yubin/src"
)

type API struct {
	yubin  *yubin.Yubin
	engine *gin.Engine
}

func New(yubin *yubin.Yubin) *API {
	api := &API{
		yubin:  yubin,
		engine: gin.New(),
	}
	api.InitRoutes()
	return api
}

func (api *API) InitRoutes() {
	pubsub.NewController(api.yubin).InitRoutes(api.engine.Group("/pubsub"))
	templates.NewController(api.yubin).InitRoutes(api.engine.Group("/templates"))
	users.NewController(api.yubin).InitRoutes(api.engine.Group("/users"))
	sources.NewController(api.yubin).InitRoutes(api.engine.Group("/sources"))
}

func (api *API) Run(addr ...string) {
	if err := api.engine.Run(addr...); err != nil {
		log.Println("failed to run http server:", err)
	}
}
