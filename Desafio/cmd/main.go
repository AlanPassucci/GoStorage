package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"desafio/cmd/router"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	db, err := sql.Open(os.Getenv("DATA_BASE"), os.Getenv("MYSQL_DATA_SOURCE"))
	if err != nil {
		panic(err)
	}

	server := gin.Default()

	router.NewRouter(server, db).MapRoutes()

	server.Run()

}
