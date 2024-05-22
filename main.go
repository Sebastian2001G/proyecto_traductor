package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	app := fiber.New()
	app.Use(cors.New()) //Se activan los cors para que se procesen las peticiones

	app.Static("/", "./client/dist") //Se le asigna la ruta para el front

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"data": "usuarios desde el backend",
		})
	})

	app.Listen(":" + port)                                    //Se inicia el puerto
	fmt.Println("Servidor ejecutandose en el puerto " + port) //Se indica que el servidor se esta ejecutando
}
