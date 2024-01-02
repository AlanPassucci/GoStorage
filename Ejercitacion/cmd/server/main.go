package main

import (
	"database/sql"
	"gostorage/cmd/server/handler"
	"gostorage/internal/product"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	mysqlDataSource := os.Getenv("MYSQL_DATA_SOURCE")

	db, err := sql.Open("mysql", mysqlDataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	repository := product.NewRepository(db)
	service := product.NewService(repository)
	productHandler := handler.NewProductHandler(service)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	products := r.Group("/products")
	{
		products.GET("/", productHandler.GetAll())
		products.GET("/:id", productHandler.GetByID())
		products.POST("/", productHandler.Create())
		products.PATCH("/:id", productHandler.Update())
		products.DELETE("/:id", productHandler.Delete())
	}

	r.Run(":8080")
}
