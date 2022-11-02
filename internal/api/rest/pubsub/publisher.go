package pubsub

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer"
	"smtp-client/pkg/data/crud"
	"time"
)

type Publisher struct {
	mailer *mailer.Mailer
}

func (p *Publisher) InitRoutes(group *gin.RouterGroup) {
	group.POST("/publish", p.publish)
	group.GET("/:id", p.getPublication)
	group.GET("/:id/report", p.getReport)
	group.GET("/:id/report/:user_id", p.getPersonalReport)
	group.GET("/reports", p.getReports)
	group.GET("/", p.getPublications)
}

func (p *Publisher) getPublication(ctx *gin.Context) {
	id := ctx.Param("id")
	pub, ok, err := p.mailer.Publications().Get(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, publication2dto(pub))
}

func (p *Publisher) getReport(ctx *gin.Context) {
	id := ctx.Param("id")
	report, ok, err := p.mailer.Reports().Get(id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, report2dto(report))
}

func (p *Publisher) getReports(ctx *gin.Context) {
	query := util.ParseCRUDQuery(ctx)
	reports, err := p.mailer.Reports().Query(query)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var dtos []ReportDto
	for _, report := range reports {
		dtos = append(dtos, report2dto(report))
	}

	dtos = util.ValidateEmptySlice(dtos)

	ctx.JSON(http.StatusOK, dtos)
}

func (p *Publisher) getPersonalReport(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.Param("user_id")
	report, ok, err := p.mailer.PersonalReports().Get(crud.PairKey[string, string]{
		Item1: id,
		Item2: userID,
	})
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !ok {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, personalReport2dto(report))
}

func (p *Publisher) getPublications(ctx *gin.Context) {
	query := util.ParseCRUDQuery(ctx)
	pubs, err := p.mailer.Publications().Query(query)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	var dtos []PublicationDto
	for _, pub := range pubs {
		dtos = append(dtos, publication2dto(pub))
	}
	dtos = util.ValidateEmptySlice(dtos)
	ctx.JSON(http.StatusOK, dtos)
}

func (p *Publisher) publish(ctx *gin.Context) {
	var publication PublicationDto
	if err := ctx.BindJSON(&publication); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var at *time.Time
	if publication.At != nil {
		at = new(time.Time)
		*at = time.Unix(*publication.At, 0)
	}

	id, err := p.mailer.Publish(mailer.Use(mailer.PublishOptions{
		SendOptions: mailer.SendOptions{
			Topics:   publication.Topics,
			Users:    publication.Users,
			SourceID: publication.Source,
			Template: publication.Template,
		},
		At:   at,
		Meta: publication.Meta,
	}))
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte(id))
}
