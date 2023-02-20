package rest

import (
	"github.com/gin-gonic/gin"
	"log"
	"smtp-client/internal/api/rest/pubsub"
	"smtp-client/internal/api/rest/sources"
	"smtp-client/internal/api/rest/templates"
	"smtp-client/internal/api/rest/users"
	"smtp-client/internal/yubin"
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
