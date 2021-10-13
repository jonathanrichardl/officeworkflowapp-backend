package main

import (
	"clean/internal/controller"
	"clean/internal/infrastructure/repository"
	"clean/internal/usecase/orders"
	"clean/internal/usecase/requirements"
	"clean/pkg/logger"
	"database/sql"
	"os"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	logger := logger.NewLogger()
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL")+"?parseTime=true")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// db, err := sql.Open("mysql", "root:123jonathan12310w0300!!!@tcp(localhost:3306)/testers?parseTime=true")
	if err != nil {
		panic(err)
	}
	orderRepo := repository.NewOrdersMySQL(db)
	requirementRepo := repository.NewRequirementsMySQL(db)
	orderService := orders.NewService(orderRepo)
	requirementService := requirements.NewService(requirementRepo)
	c := controller.NewController(orderService, requirementService, logger)
	c.RegisterHandler()
	c.Start()

}
