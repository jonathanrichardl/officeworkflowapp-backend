package main

import (
	"database/sql"
	"order-validation-v2/internal/controller"
	"order-validation-v2/internal/infrastructure/repository"
	"order-validation-v2/internal/usecase/orders"
	"order-validation-v2/internal/usecase/requirements"
	"order-validation-v2/internal/usecase/submissions"
	"order-validation-v2/internal/usecase/tasks"
	"order-validation-v2/internal/usecase/user"
	"order-validation-v2/pkg/logger"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	logger := logger.NewLogger()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	taskRepo := repository.NewTaskPSQL(db)
	submissionRepo := repository.NewSubmissionPSQL(db)
	orderRepo := repository.NewOrdersPSQL(db)
	requirementRepo := repository.NewRequirementsPSQL(db)
	userRepo := repository.NewUserPSQL(db)
	orderService := orders.NewService(orderRepo)
	requirementService := requirements.NewService(requirementRepo)
	userService := user.NewService(userRepo)
	taskService := tasks.NewService(taskRepo)
	submissionService := submissions.NewService(submissionRepo)
	c := controller.NewController(orderService, userService, requirementService,
		taskService, submissionService, logger)
	c.RegisterHandler()
	c.Start()

}
