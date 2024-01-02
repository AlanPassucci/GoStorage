package handler

import (
	"desafio/internal/domain"
	"desafio/internal/invoices"

	"github.com/gin-gonic/gin"
)

type Invoices struct {
	s invoices.Service
}

func NewHandlerInvoices(s invoices.Service) *Invoices {
	return &Invoices{s}
}

func (i *Invoices) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		invoices, err := i.s.ReadAll(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, invoices)
	}
}

func (i *Invoices) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		invoices := domain.Invoice{}
		err := ctx.ShouldBindJSON(&invoices)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = i.s.Create(ctx, &invoices)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"data": invoices})
	}
}

func (i *Invoices) PostManyFromJSON() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		invoices, err := i.s.CreateManyFromJSON(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"data": invoices})
	}
}

func (i *Invoices) UpdateTotals() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := i.s.UpdateTotals(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"data": "Totals updated successfully"})
	}
}
