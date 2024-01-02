package handler

import (
	"desafio/internal/customers"
	"desafio/internal/domain"

	"github.com/gin-gonic/gin"
)

type Customers struct {
	s customers.Service
}

func NewHandlerCustomers(s customers.Service) *Customers {
	return &Customers{s}
}

func (c *Customers) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customers, err := c.s.ReadAll(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, customers)
	}
}

func (c *Customers) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customer := domain.Customer{}
		err := ctx.ShouldBindJSON(&customer)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = c.s.Create(ctx, &customer)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(201, gin.H{"data": customer})
	}
}

func (c *Customers) PostManyFromJSON() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customers, err := c.s.CreateManyFromJSON(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, gin.H{"data": customers})
	}
}

func (c *Customers) GetTotalsGroupedByCondition() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		totalGrouped, err := c.s.GetTotalsGroupedByCondition(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"data": totalGrouped})
	}
}

func (c *Customers) GetActivesWhoSpentTheMost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		activesWhoSpentTheMost, err := c.s.GetActivesWhoSpentTheMost(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"data": activesWhoSpentTheMost})
	}
}
