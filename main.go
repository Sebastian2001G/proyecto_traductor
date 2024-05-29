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
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://proyectotraductor62:t08kHYaPpvUmA03R@traductor.tabvz9q.mongodb.net/?retryWrites=true&w=majority&appName=Traductor").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	app.Use(cors.New()) //Se activan los cors para que se procesen las peticiones

	app.Static("/", "./client/dist") //Se le asigna la ruta para el front

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

	app.Get("/buscar", func(c *fiber.Ctx) error {
		palabra := c.Query("palabra")
		if palabra == "" {
			return c.Status(400).SendString("Missing query parameter: palabra")
		}

		coll := client.Database("gomongodb").Collection("words")
		filter := bson.M{"WordText": palabra}

		var result models.Word
		err := coll.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{
					"statusCode": 404,
					"error":      "No se encontró la palabra",
				})
			}
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(result)
	})

	app.Get("/sugerencias", func(c *fiber.Ctx) error {
		palabra := c.Query("palabra")
		if palabra == "" {
			return c.Status(400).SendString("Missing query parameter: palabra")
		}

		primeraLetra := string([]rune(palabra)[0])
		coll := client.Database("gomongodb").Collection("words")
		filter := bson.M{"WordText": bson.M{"$regex": "^" + primeraLetra}}
		opts := options.Find().SetLimit(5)

		cursor, err := coll.Find(context.TODO(), filter, opts)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer cursor.Close(context.TODO())

		var results []models.Word
		if err = cursor.All(context.TODO(), &results); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if len(results) == 0 {
			return c.Status(404).SendString("No se encontraron palabras")
		}

		return c.JSON(results)
	})

	app.Listen(":" + port)                                    //Se inicia el puerto
	fmt.Println("Servidor ejecutandose en el puerto " + port) //Se indica que el servidor se esta ejecutando
}
