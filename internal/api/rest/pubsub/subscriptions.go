package pubsub

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer"
)

type Subscriptions struct {
	mailer *mailer.Mailer
}

func (s *Subscriptions) InitRoutes(group *gin.RouterGroup) {
	group.GET("/:id", s.getSubscriptions)
	group.GET("/:id/:topic", s.getSubscription)
	group.POST("/:id/:topic", s.subscribe)
	group.PUT("/:id/:topic", s.update)
	group.DELETE("/:id/:topic", s.unsubscribe)
	group.DELETE("/:id", s.unsubscribeAll)
}

func (s *Subscriptions) getSubscription(ctx *gin.Context) {
	id := ctx.Param("id")
	topic := ctx.Param("topic")
	subscription, ok, err := s.mailer.Subscriptions().Subscription(id, topic)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, subscription2dto(subscription))
}

func (s *Subscriptions) getSubscriptions(ctx *gin.Context) {
	id := ctx.Param("id")
	subscriptions, err := s.mailer.Subscriptions().Subscriptions(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var dtos []SubscriptionDto
	for _, subscription := range subscriptions {
		dtos = append(dtos, subscription2dto(subscription))
	}

	dtos = util.ValidateEmptySlice(dtos)

	ctx.JSON(http.StatusOK, dtos)
}

func (s *Subscriptions) subscribe(ctx *gin.Context) {
	topic := ctx.Param("topic")
	id := ctx.Param("id")
	subscription, err := s.mailer.Subscriptions().Subscribe(id, topic)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, subscription2dto(subscription))
}

func (s *Subscriptions) unsubscribe(ctx *gin.Context) {
	topic := ctx.Param("topic")
	id := ctx.Param("id")
	err := s.mailer.Subscriptions().Unsubscribe(id, topic)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (s *Subscriptions) unsubscribeAll(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.mailer.Subscriptions().UnsubscribeAll(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (s *Subscriptions) update(ctx *gin.Context) {
	topic := ctx.Param("topic")
	id := ctx.Param("id")

	var updateDto UpdateSubscriptionDto
	if err := ctx.BindJSON(&updateDto); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	sub, ok, err := s.mailer.Subscriptions().Subscription(id, topic)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}

	sub.Meta = updateDto.Meta
	err = s.mailer.Subscriptions().Update(sub)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, subscription2dto(sub))
}
