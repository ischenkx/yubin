package users

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"yubin/common/data"
	"yubin/lib/api/rest/util"
	yubin "yubin/src"
	"yubin/src/user"
)

type Controller struct {
	mailer *yubin.Yubin
}

func NewController(mailer *yubin.Yubin) *Controller {
	return &Controller{mailer: mailer}
}

func (c *Controller) get(ctx *gin.Context) {
	users, err := util.KvRangeQuery[user.User](ctx, c.mailer.Users(), util.ParseRangeQuery(ctx))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var dtos []UserDto
	for _, u := range users {
		dtos = append(dtos, user2dto(u))
	}
	dtos = util.ValidateEmptySlice(dtos)
	ctx.JSON(http.StatusOK, dtos)
}

func (c *Controller) getByID(ctx *gin.Context) {
	id := ctx.Param("id")

	u, err := c.mailer.Users().Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	ctx.JSON(http.StatusOK, user2dto(u))
}

func (c *Controller) delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.mailer.Subscriptions().UnsubscribeAll(ctx, id); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := c.mailer.Users().Delete(ctx, id)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
}

func (c *Controller) update(ctx *gin.Context) {
	id := ctx.Param("id")
	var dto UpdateDto
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	u, err := c.mailer.Users().Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	if dto.Name != nil {
		u.Name = *dto.Name
	}
	if dto.Surname != nil {
		u.Surname = *dto.Surname
	}
	if dto.Email != nil {
		u.Email = *dto.Email
	}
	if dto.Meta != nil {
		u.Meta = dto.Meta
	}
	err = c.mailer.Users().Set(ctx, u.ID, u)
	if err != nil {
		log.Println("failed to update a user:", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, user2dto(u))
}

func (c *Controller) create(ctx *gin.Context) {
	var dto CreateDto

	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	u := user.User{
		ID:      uuid.New().String(),
		Email:   dto.Email,
		Name:    dto.Name,
		Surname: dto.Surname,
		Meta:    dto.Meta,
	}

	err := c.mailer.Users().Set(ctx, u.ID, u)
	if err != nil {
		ctx.Status(http.StatusNotModified)
		return
	}
	ctx.JSON(http.StatusOK, user2dto(u))
}

func (c *Controller) InitRoutes(group *gin.RouterGroup) {
	group.GET("/", c.get)
	group.GET("/:id", c.getByID)
	group.DELETE("/:id", c.delete)
	group.PUT("/", c.update)
	group.POST("/", c.create)
}
