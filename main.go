package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type QnA []struct {
	Question string   `json:"Question"`
	Choices  []string `json:"Choices"`
	Answer   int      `json:"Answer"`
}

type Cards struct {
	Proj []struct {
		Img   string `json:"img"`
		Title string `json:"title"`
		Txt   string `json:"txt"`
		Link  string `json:"link"`
	} `json:"proj"`
	Exp []struct {
		Img   string `json:"img"`
		Title string `json:"title"`
		Txt   string `json:"txt"`
		Link  string `json:"link"`
	} `json:"exp"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("This is the frontpage")
	})

	app.Get("/sup", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Sup ✋",
		})
	})

	app.Get("/gamesshow", func(c *fiber.Ctx) error {
		jsonFile, err := os.Open("./resources/QnA.json")
		if err != nil {
			return c.JSON(fiber.Map{"Message": "Unable to open file"})
		}
		defer jsonFile.Close()

		jsonData, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return c.JSON(fiber.Map{"Message": "Unable to read file"})
		}
		data := QnA{}
		err2 := json.Unmarshal(jsonData, &data)
		if err2 != nil {
			return c.JSON(fiber.Map{"Message": "Unable to read json data"})
		}
		return c.JSON(data)

	})

	app.Listen(":8000")

}
