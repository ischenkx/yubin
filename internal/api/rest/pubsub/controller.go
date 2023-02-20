package pubsub

import (
	"github.com/gin-gonic/gin"
	"smtp-client/internal/yubin"
)

type Controller struct {
	mailer *yubin.Yubin
}

func NewController(mailer *yubin.Yubin) *Controller {
	return &Controller{mailer: mailer}
}

func (c *Controller) InitRoutes(group *gin.RouterGroup) {
	topics := &Topics{mailer: c.mailer}
	publisher := &Publisher{mailer: c.mailer}
	subscriptions := &Subscriptions{mailer: c.mailer}

	topics.InitRoutes(group.Group("/topics"))
	publisher.InitRoutes(group.Group("/publisher"))
	subscriptions.InitRoutes(group.Group("/subscriptions"))
}
