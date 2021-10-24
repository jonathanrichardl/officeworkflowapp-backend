package main

import (
	"database/sql"
	"order-validation-v2/internal/controller"
	"order-validation-v2/internal/infrastructure/repository"
	"order-validation-v2/internal/usecase/orders"
	"order-validation-v2/internal/usecase/requirements"
	"order-validation-v2/internal/usecase/user"
	"order-validation-v2/pkg/logger"
	"os"

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
	orderRepo := repository.NewOrdersPSQL(db)
	requirementRepo := repository.NewRequirementsPSQL(db)
	userRepo := repository.NewUserPSQL(db)
	orderService := orders.NewService(orderRepo)
	requirementService := requirements.NewService(requirementRepo)
	userService := user.NewService(userRepo)

	c := controller.NewController(orderService, userService, requirementService, logger)
	c.RegisterHandler()
	c.Start()

}
