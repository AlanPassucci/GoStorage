package handler

import (
	"errors"
	"gostorage/internal/domain"
	"gostorage/internal/product"
	"gostorage/pkg/web"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	sv product.Service
}

func NewProductHandler(sv product.Service) *ProductHandler {
	return &ProductHandler{sv}
}

func (h *ProductHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// process
		// - get all products
		products, err := h.sv.GetAll(ctx)
		if err != nil {
			web.Failure(ctx, http.StatusInternalServerError, err)
			return
		}

		// response
		web.Success(ctx, http.StatusOK, products)
	}
}

func (h *ProductHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// - get id from path
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}

		// process
		// - get product by id
		p, err := h.sv.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrServiceInvalidProductID):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceProductNotFound):
				web.Failure(ctx, http.StatusNotFound, err)
			default:
				web.Failure(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		// response
		web.Success(ctx, http.StatusOK, p)
	}
}

func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// - get product from body
		var p domain.Product
		if err := ctx.ShouldBindJSON(&p); err != nil {
			web.Failure(ctx, http.StatusUnprocessableEntity, err)
			return
		}

		// process
		// - create product
		p, err := h.sv.Create(ctx, &p)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrServiceInvalidProductName):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductQuantity):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductCodeValue):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductExpiration):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductPrice):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrorServiceAlreadyExistsCodeValue):
				web.Failure(ctx, http.StatusConflict, err)
			default:
				web.Failure(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		// response
		web.Success(ctx, http.StatusCreated, p)
	}
}

func (h *ProductHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// - get id from path
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}

		// process
		// - get product from body
		originalProduct, err := h.sv.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrServiceInvalidProductID):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceProductNotFound):
				web.Failure(ctx, http.StatusNotFound, err)
			default:
				web.Failure(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		// - get product from body
		productToUpdate := domain.Product{
			Id:          originalProduct.Id,
			Name:        originalProduct.Name,
			Quantity:    originalProduct.Quantity,
			CodeValue:   originalProduct.CodeValue,
			IsPublished: originalProduct.IsPublished,
			Expiration:  originalProduct.Expiration,
			Price:       originalProduct.Price,
		}
		if err := ctx.ShouldBindJSON(&productToUpdate); err != nil {
			web.Failure(ctx, http.StatusUnprocessableEntity, err)
			return
		}

		if originalProduct == productToUpdate {
			web.Success(ctx, http.StatusNoContent, nil)
			return
		}

		// - update product
		if err := h.sv.Update(ctx, &productToUpdate); err != nil {
			switch {
			case errors.Is(err, product.ErrServiceInvalidProductID):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductName):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductQuantity):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductCodeValue):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductExpiration):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceInvalidProductPrice):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrorServiceAlreadyExistsCodeValue):
				web.Failure(ctx, http.StatusConflict, err)
			default:
				web.Failure(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		// response
		web.Success(ctx, http.StatusOK, productToUpdate)
	}
}

func (h *ProductHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// - get id from path
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}

		// process
		// - delete product by id
		if err := h.sv.Delete(ctx, id); err != nil {
			switch {
			case errors.Is(err, product.ErrServiceInvalidProductID):
				web.Failure(ctx, http.StatusBadRequest, err)
			case errors.Is(err, product.ErrServiceProductNotFound):
				web.Failure(ctx, http.StatusNotFound, err)
			default:
				web.Failure(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		// response
		web.Success(ctx, http.StatusNoContent, nil)
	}
}
