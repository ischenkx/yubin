package users

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer"
	"smtp-client/internal/mailer/user"
)

type Controller struct {
	mailer *mailer.Mailer
}

func NewController(mailer *mailer.Mailer) *Controller {
	return &Controller{mailer: mailer}
}

func (c *Controller) get(ctx *gin.Context) {
	query := util.ParseCRUDQuery(ctx)
	users, err := c.mailer.Users().Query(query)
	if err != nil {
		log.Println("failed to get users:", err)
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

	u, ok, err := c.mailer.Users().Get(id)
	if err != nil {
		log.Println("failed to get a user by id:", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, user2dto(u))
}

func (c *Controller) delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.mailer.Subscriptions().UnsubscribeAll(id); err != nil {
		log.Println("failed to unsubscribe a user from all topics:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := c.mailer.Users().Delete(id)
	if err != nil {
		log.Println("failed to delete a user:", err)
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
	u, ok, err := c.mailer.Users().Get(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
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
	err = c.mailer.Users().Update(u)
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
	res, err := c.mailer.Users().Create(user.User{
		Email:   dto.Email,
		Name:    dto.Name,
		Surname: dto.Surname,
		Meta:    dto.Meta,
	})
	if err != nil {
		ctx.Status(http.StatusNotModified)
		return
	}
	ctx.JSON(http.StatusOK, user2dto(res))
}

func (c *Controller) InitRoutes(group *gin.RouterGroup) {
	group.GET("/", c.get)
	group.GET("/:id", c.getByID)
	group.DELETE("/:id", c.delete)
	group.PUT("/", c.update)
	group.POST("/", c.create)
}
