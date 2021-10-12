package main

import (
	"clean/internal/controller"
	"clean/internal/infrastructure/repository"
	"clean/internal/usecase/orders"
	"clean/internal/usecase/requirements"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123jonathan123100300!!!@tcp(127.0.0.1:3306)/testers")
	if err != nil {
		panic(err)
	}
	orderRepo := repository.NewOrdersMySQL(db)
	requirementRepo := repository.NewRequirementsMySQL(db)
	orderService := orders.NewService(orderRepo)
	requirementService := requirements.NewService(requirementRepo)
	c := controller.NewController(orderService, requirementService)
	c.RegisterHandler()
	c.Start()

}
