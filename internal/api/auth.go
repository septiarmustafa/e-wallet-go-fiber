package api

import (
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"

	"github.com/gofiber/fiber/v2"
)

type AuthApi struct {
	userService domain.UserService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMid fiber.Handler) {
	h := AuthApi{
		userService: userService,
	}
	app.Post("token/generate", h.GenerateToken)
	app.Get("token/validate", authMid, h.ValidateToken)
	app.Post("user/register", h.RegisterUser)
	app.Post("user/validate-otp", h.ValidateOTP)
}

func (a AuthApi) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	token, err := a.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}
	return ctx.Status(200).JSON(token)
}

func (a AuthApi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")
	return ctx.Status(200).JSON(user)
}

func (a AuthApi) RegisterUser(ctx *fiber.Ctx) error {
	var req dto.UserRegisterReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	res, err := a.userService.Register(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}
	return ctx.Status(200).JSON(res)
}

func (a AuthApi) ValidateOTP(ctx *fiber.Ctx) error {
	var req dto.ValidateOtpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	err := a.userService.ValidateOTP(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}
	return ctx.SendStatus(200)
}
