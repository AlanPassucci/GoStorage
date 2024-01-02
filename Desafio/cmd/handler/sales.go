package handler

import (
	"desafio/internal/domain"
	"desafio/internal/sales"

	"github.com/gin-gonic/gin"
)

type Sales struct {
	s sales.Service
}

func NewHandlerSales(s sales.Service) *Sales {
	return &Sales{s}
}

func (s *Sales) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		invoices, err := s.s.ReadAll(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, invoices)
	}
}

func (s *Sales) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sale := domain.Sale{}
		err := ctx.ShouldBindJSON(&sale)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = s.s.Create(ctx, &sale)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"data": sale})
	}
}

func (s *Sales) PostManyFromJSON() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sales, err := s.s.CreateManyFromJSON(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"data": sales})
	}
}
