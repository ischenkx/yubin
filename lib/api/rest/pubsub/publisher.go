package pubsub

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yubin/common/data"
	"yubin/common/data/record"
	"yubin/lib/api/rest/util"
	"yubin/lib/piper"
	yubin "yubin/src"
	"yubin/src/publication"
)

type Publisher struct {
	mailer *yubin.Yubin
	piper  *piper.Piper
}

func (p *Publisher) InitRoutes(group *gin.RouterGroup) {
	group.POST("/publish", p.publish)
	group.GET("/:id", p.getPublication)
	group.GET("/:id/reports", p.getPublicationReports)
	group.GET("/:id/reports/:user_id", p.getUserReport)
	group.GET("/reports", p.getReports)
	group.GET("/", p.getPublications)
}

func (p *Publisher) getPublicationReports(ctx *gin.Context) {
	id := ctx.Param("id")
	set := p.mailer.Reports().
		Filter(record.E{"publication", id})

	reports, err := util.RecordsQuery(ctx, set, "publication", util.ParseRangeQuery(ctx))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}

	var dtos []ReportDto
	for _, report := range reports {
		dtos = append(dtos, recordReport2dto(report))
	}

	dtos = util.ValidateEmptySlice(dtos)

	ctx.JSON(http.StatusOK, dtos)
}

func (p *Publisher) getPublication(ctx *gin.Context) {
	id := ctx.Param("id")
	pub, err := p.mailer.Publications().Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, publication2dto(pub))
}

func (p *Publisher) getReports(ctx *gin.Context) {
	reports, err := util.RecordsQuery(ctx, p.mailer.Reports(), "publication", util.ParseRangeQuery(ctx))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var dtos []ReportDto
	for _, report := range reports {
		dtos = append(dtos, recordReport2dto(report))
	}

	dtos = util.ValidateEmptySlice(dtos)

	ctx.JSON(http.StatusOK, dtos)
}

func (p *Publisher) getUserReport(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.Param("user_id")
	report, err := p.mailer.Reports().
		Filter(record.E{"publication", id}, record.E{"user", userID}).
		Cursor().
		Iter().
		Next(ctx)

	if err != nil {
		if err == data.NotFoundErr {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
		return
	}
	ctx.JSON(http.StatusOK, recordReport2dto(report))
}

func (p *Publisher) getPublications(ctx *gin.Context) {
	pubs, err := util.KvRangeQuery[publication.Publication](ctx, p.mailer.Publications(), util.ParseRangeQuery(ctx))
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

	options := []yubin.PublishOption{
		yubin.ToTopics(publication.Topics...),
		yubin.ToUsers(publication.Users...),
		yubin.WithSource(publication.Source),
		yubin.WithTemplate(publication.Template),
	}

	for name, val := range publication.Properties {
		options = append(options, yubin.WithProperty(name, val))
	}

	pub, err := p.mailer.New(ctx, options...)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if err := p.piper.Schedule(ctx, pub.ID); err != nil {
		log.Println("failed to schedule:", err)
		ctx.Status(http.StatusInternalServerError)
	}

	ctx.Data(http.StatusOK, gin.MIMEPlain, []byte(pub.ID))
}
