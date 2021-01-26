package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/atoato88/url-shortener/pkg/conf"
	"github.com/atoato88/url-shortener/pkg/data"
	"github.com/atoato88/url-shortener/pkg/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var urls []data.URLentry
var co = conf.Conf

func saveFile(urls map[string]data.URLentry) error {
	file, err := json.MarshalIndent(urls, "", " ")
	err = ioutil.WriteFile(util.FilePath, file, 0644)

	return err
}

func loadDB() error {
	log.Printf("read db from: %s", util.FilePath)
	content, _ := ioutil.ReadFile(util.FilePath)
	err := json.Unmarshal(content, &urls)
	if err != nil {
		log.Printf("%v", err)
	}
	if co.Debug {
		log.Printf("initialized db with: %v", urls)
	}
	return nil
}

func renderString(c *fiber.Ctx) error {
	now := time.Now()
	return c.SendString(strings.Join([]string{fmt.Sprintf("rendered at %s", now)}, "\n"))
}

func main() {
	co.Debug = true

	port := flag.Int("port", co.Port, "port number listening for")
	flag.Parse()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/:id", handleRedirect)

	log.Println("start server")
	log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
}

func handleRedirect(c *fiber.Ctx) error {
	loadDB()
	id := c.Params("id")

	for _, e := range urls {
		if e.Id == id {
			log.Println(fmt.Sprintf("redirect to \"%s\" requested from %s", e.Url, c.IP()))
			c.Set("Cache-Control", "no-store")
			return c.Redirect(e.Url, 301)
		}
	}

	return c.SendString(fmt.Sprintf("no URL entry for id: \"%s\"", string(id)))
}
