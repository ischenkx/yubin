package sources

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yubin/common/data"
	"yubin/lib/api/rest/util"
	yubin "yubin/src"
)

type Controller struct {
	mailer *yubin.Yubin
}

func NewController(mailer *yubin.Yubin) *Controller {
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
	err := c.mailer.Sources().Set(ctx, dto.Name, dto2source(dto))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) deleteSource(ctx *gin.Context) {
	id := ctx.Param("name")
	err := c.mailer.Sources().Delete(ctx, id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) getSource(ctx *gin.Context) {
	id := ctx.Param("name")
	source, err := c.mailer.Sources().Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	ctx.JSON(http.StatusOK, source2dto(source))
}

func (c *Controller) getSources(ctx *gin.Context) {
	sources, err := util.KvRangeQuery[yubin.NamedSource](ctx, c.mailer.Sources(), util.ParseRangeQuery(ctx))
	if err != nil {
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

	s, err := c.mailer.Sources().Get(ctx, dto.Name)
	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
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

	err = c.mailer.Sources().Set(ctx, s.Name, s)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
