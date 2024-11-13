package main

import (
	"e-wallet/dto"
	"e-wallet/internal/api"
	"e-wallet/internal/component"
	"e-wallet/internal/config"
	"e-wallet/internal/middleware"
	"e-wallet/internal/repository"
	"e-wallet/internal/service"
	"e-wallet/internal/sse"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := component.GetDatabaseConnection(cnf)
	cacheConnection := repository.NewRedisClient(cnf)

	hub := &dto.Hub{
		NotificationChannel: map[int64]chan dto.NotificationData{},
	}

	userRespository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)

	emailService := service.NewEmail(cnf)
	userService := service.NewUser(userRespository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationRepository, hub)
	notificationService := service.NewNotification(notificationRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()
	api.NewAuth(app, userService, authMid)
	api.NewTransfer(app, authMid, transactionService)
	api.NewNotification(app, authMid, notificationService)
	sse.NewNotification(app, authMid, hub)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
