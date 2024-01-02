package handler

import (
	"desafio/internal/domain"
	"desafio/internal/products"

	"github.com/gin-gonic/gin"
)

type Products struct {
	s products.Service
}

func NewHandlerProducts(s products.Service) *Products {
	return &Products{s}
}

func (p *Products) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := p.s.ReadAll(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, products)
	}
}

func (p *Products) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products := domain.Product{}
		err := ctx.ShouldBindJSON(&products)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		err = p.s.Create(ctx, &products)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, gin.H{"data": products})
	}
}

func (p *Products) PostManyFromJSON() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := p.s.CreateManyFromJSON(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, gin.H{"data": products})
	}
}

func (p *Products) GetQtySaledGroupedByDescription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		qtySaledGrouped, err := p.s.GetQtySaledGroupedByDescription(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"data": qtySaledGrouped})
	}
}
