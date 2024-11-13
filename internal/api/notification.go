package api

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
)

type notificationApi struct {
	notificationService domain.NotificationService
}

func NewNotification(app *fiber.App, authMid fiber.Handler, notificationService domain.NotificationService) {
	h := notificationApi{
		notificationService: notificationService,
	}

	app.Get("/notifications", authMid, h.GetUsersNotification)
}

func (n notificationApi) GetUsersNotification(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 15 * time.Second)
	defer cancel()

	// use local user
	user := ctx.Locals("x-user").(dto.UserData)

	notification, err := n.notificationService.FindByUser(c, user.ID)

	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{ Message: err.Error()})
	}

	return ctx.Status(200).JSON(notification)

}
