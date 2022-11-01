package util

import (
	"github.com/gin-gonic/gin"
	"smtp-client/pkg/data/crud"
)

type queryDto struct {
	Offset *int `json:"offset" form:"offset"`
	Limit  *int `json:"limit" form:"limit"`
}

func ParseCRUDQuery(ctx *gin.Context) *crud.Query {
	var query queryDto
	if err := ctx.BindQuery(&query); err != nil {
		return nil
	}
	return &crud.Query{
		Offset: query.Offset,
		Limit:  query.Limit,
	}
}
