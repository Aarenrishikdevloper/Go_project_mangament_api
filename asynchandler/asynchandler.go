package asynchandler

import "github.com/gofiber/fiber/v3"

type AsyncControllerType func(c fiber.Ctx) error

func AsyncHandler(controller AsyncControllerType) AsyncControllerType {
	return func(c fiber.Ctx) error {
		if err := controller(c); err != nil {
			return err
		}
		return nil
	}

}
