package helper

import (
	"github.com/gofiber/fiber/v2"

	"go_module5/app/model"
)

const LocalsAuthUser = "authUser"

func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	user, ok := c.Locals(LocalsAuthUser).(model.AuthUser)
	return user, ok
}
