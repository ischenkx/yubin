package templates

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yubin/common/data"
	"yubin/lib/api/rest/util"
	yubin "yubin/src"
	"yubin/src/template"
)

type Controller struct {
	mailer *yubin.Yubin
}

func NewController(mailer *yubin.Yubin) *Controller {
	return &Controller{mailer: mailer}
}

func (c *Controller) InitRoutes(group *gin.RouterGroup) {
	group.GET("/", c.getTemplates)
	group.GET("/:id", c.getTemplate)
	group.DELETE("/:id", c.deleteTemplate)
	group.POST("/", c.createTemplate)
	group.PUT("/", c.updateTemplate)
}

func (c *Controller) createTemplate(ctx *gin.Context) {
	var dto TemplateDto
	if err := ctx.BindJSON(&dto); err != nil {
		log.Println(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	err := c.mailer.Templates().Set(ctx, dto.Name(), dto)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) deleteTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.mailer.Templates().Delete(ctx, id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) getTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	t, err := c.mailer.Templates().Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	ctx.JSON(http.StatusOK, template2dto(t))
}

func (c *Controller) getTemplates(ctx *gin.Context) {
	templates, err := util.KvRangeQuery[template.Template](ctx, c.mailer.Templates(), util.ParseRangeQuery(ctx))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	output := make([]TemplateDto, 0, len(templates))
	for _, t := range templates {
		output = append(output, template2dto(t))
	}

	ctx.JSON(http.StatusOK, output)
}

func (c *Controller) updateTemplate(ctx *gin.Context) {
	var dto UpdateTemplateDto
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	t, err := c.mailer.Templates().Get(ctx, dto.TemplateName)

	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	model := template2dto(t)

	if dto.Sub != nil {
		model.Sub = *dto.Sub
	}

	if dto.Data != nil {
		model.Data = *dto.Data
	}

	if dto.Meta != nil {
		model.MetaData = *dto.Meta
	}

	err = c.mailer.Templates().Set(ctx, model.TemplateName, model)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
