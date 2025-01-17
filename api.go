package main

import "github.com/gofiber/fiber/v3"

// fetch data from db
func DataHandler(c fiber.Ctx) error {
    return c.SendString("data: not implemented")
}

// upload new data to db
func UploadHandler(c fiber.Ctx) error {
    data := new(Datapoint)

    if err := c.Bind().Body(data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid data",
        })
    }

    if err := WriteDatapoint(data); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to write data to database",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Data successfully saved",
    })
}
