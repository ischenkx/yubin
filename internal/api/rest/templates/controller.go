package templates

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
	_, err := c.mailer.Templates().Create(dto)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) deleteTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.mailer.Templates().Delete(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) getTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	t, ok, err := c.mailer.Templates().Get(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, template2dto(t))
}

func (c *Controller) getTemplates(ctx *gin.Context) {
	query := util.ParseCRUDQuery(ctx)
	templates, err := c.mailer.Templates().Query(query)
	if err != nil {
		log.Println(err)
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

	t, ok, err := c.mailer.Templates().Get(dto.TemplateName)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
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

	err = c.mailer.Templates().Update(model)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
}
