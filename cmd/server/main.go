package main

import (
	"database/sql"

	"github.com/davidop97/apiGo/cmd/server/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// swagger documentation
// @title API GO
// @version 1.0
// @description This API manage many products of any company.
// @host localhost:8080/api/v1
func main() {
	// NO MODIFICAR
	db, err := sql.Open("mysql", "mysql_apigo_user:MySql_ApiGo#97@/mysqlapigo")
	if err != nil {
		panic(err)
	}

	eng := gin.Default()

	router := routes.NewRouter(eng, db)
	router.MapRoutes()

	if err := eng.Run(); err != nil {
		panic(err)
	}
}
