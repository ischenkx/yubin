package pubsub

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer"
)

type Topics struct {
	mailer *mailer.Mailer
}

func (t *Topics) InitRoutes(group *gin.RouterGroup) {
	group.GET("/", t.getTopics)
	group.GET("/:topic/subscribers", t.getSubscribers)
	group.DELETE("/:topic", t.deleteTopic)
}

func (t *Topics) getSubscribers(ctx *gin.Context) {
	topic := ctx.Param("topic")
	query := util.ParseCRUDQuery(ctx)
	subscribers, err := t.mailer.Subscriptions().Subscribers(topic, query)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	subscribers = util.ValidateEmptySlice(subscribers)
	ctx.JSON(http.StatusOK, subscribers)
}

func (t *Topics) getTopics(ctx *gin.Context) {
	query := util.ParseCRUDQuery(ctx)
	topics, err := t.mailer.Subscriptions().Topics(query)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	topics = util.ValidateEmptySlice(topics)
	ctx.JSON(http.StatusOK, topics)
}

func (t *Topics) deleteTopic(ctx *gin.Context) {
	topic := ctx.Param("topic")
	err := t.mailer.Subscriptions().DeleteTopic(topic)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
