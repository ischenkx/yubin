package sources

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer"
)

type Controller struct {
	mailer *mailer.Mailer
}

func NewController(mailer *mailer.Mailer) *Controller {
	return &Controller{mailer: mailer}
}

func (c *Controller) InitRoutes(group *gin.RouterGroup) {
	group.GET("/", c.getSources)
	group.GET("/:name", c.getSource)
	group.DELETE("/:name", c.deleteSource)
	group.POST("/", c.createSource)
	group.PUT("/", c.updateSources)
}

func (c *Controller) createSource(ctx *gin.Context) {
	var dto SourceDto
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	_, err := c.mailer.Sources().Create(dto2source(dto))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) deleteSource(ctx *gin.Context) {
	id := ctx.Param("name")
	err := c.mailer.Sources().Delete(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) getSource(ctx *gin.Context) {
	id := ctx.Param("name")
	source, ok, err := c.mailer.Sources().Get(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, source2dto(source))
}

func (c *Controller) getSources(ctx *gin.Context) {
	query := util.ParseCRUDQuery(ctx)
	sources, err := c.mailer.Sources().Query(query)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	output := make([]SourceDto, 0, len(sources))
	for _, s := range sources {
		output = append(output, source2dto(s))
	}

	output = util.ValidateEmptySlice(output)

	ctx.JSON(http.StatusOK, output)
}

func (c *Controller) updateSources(ctx *gin.Context) {
	var dto UpdateSourceDto
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	s, ok, err := c.mailer.Sources().Get(dto.Name)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}

	if dto.Host != nil {
		s.Host = *dto.Host
	}
	if dto.Port != nil {
		s.Port = *dto.Port
	}
	if dto.Password != nil {
		s.Password = *dto.Password
	}
	if dto.Address != nil {
		s.Address = *dto.Address
	}

	err = c.mailer.Sources().Update(s)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
