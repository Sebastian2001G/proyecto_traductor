package main

import (
	"context"
	"fmt"
	"go-fiber-translator/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// user: proyectotraductor62 pass: t08kHYaPpvUmA03R
func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	app := fiber.New()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("MONGODB_URI is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	app.Use(cors.New()) //Se activan los cors para que se procesen las peticiones

	app.Static("/", "./client/dist") //Se le asigna la ruta para el front

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"data": "usuarios desde el backend",
		})
	})

	app.Post("/palabra", func(c *fiber.Ctx) error {
		var word models.Word

		c.BodyParser(&word)

		coll := client.Database("gomongodb").Collection("words")
		result, err := coll.InsertOne(context.TODO(), bson.D{{
			Key:   "WordText",
			Value: word.WordText,
		}, {
			Key:   "Translation",
			Value: word.Translation,
		}})

		if err != nil {
			panic(err)
		}

		return c.JSON(&fiber.Map{
			"data": result,
		})
	})

	app.Get("/palabras", func(c *fiber.Ctx) error {
		var words []models.Word
		coll := client.Database("gomongodb").Collection("words")

		results, err := coll.Find(context.TODO(), bson.M{})

		if err != nil {
			panic(err)
		}

		for results.Next(context.TODO()) {
			var word models.Word
			err := results.Decode(&word)
			if err != nil {
				panic(err)
			}
			words = append(words, word)
		}

		return c.JSON(&fiber.Map{
			"word": words,
		})
	})

	app.Listen(":" + port)                                    //Se inicia el puerto
	fmt.Println("Servidor ejecutandose en el puerto " + port) //Se indica que el servidor se esta ejecutando
}
