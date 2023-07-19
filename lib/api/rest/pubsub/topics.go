package pubsub

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yubin/lib/api/rest/util"
	yubin "yubin/src"
)

type Topics struct {
	mailer *yubin.Yubin
}

func (t *Topics) InitRoutes(group *gin.RouterGroup) {
	group.GET("/", t.getTopics)
	group.GET("/:topic/subscribers", t.getSubscribers)
	group.DELETE("/:topic", t.deleteTopic)
}

func (t *Topics) getSubscribers(ctx *gin.Context) {
	topic := ctx.Param("topic")
	subscribers, err := t.mailer.Subscriptions().Subscribers(ctx, topic)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	subscribers = util.ValidateEmptySlice(subscribers)
	ctx.JSON(http.StatusOK, subscribers)
}

func (t *Topics) getTopics(ctx *gin.Context) {
	topics, err := t.mailer.Subscriptions().Topics(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	topics = util.ValidateEmptySlice(topics)
	ctx.JSON(http.StatusOK, topics)
}

func (t *Topics) deleteTopic(ctx *gin.Context) {
	topic := ctx.Param("topic")
	err := t.mailer.Subscriptions().DeleteTopic(ctx, topic)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
