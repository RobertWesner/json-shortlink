package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	jsonFile, err := os.Open("links.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = jsonFile.Close()
	}()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	links := map[string]string{}
	err = json.Unmarshal(bytes, &links)
	if err != nil {
		log.Fatal(err)
	}

	for from, to := range links {
		r.GET(from, func(c *gin.Context) {
			c.Redirect(302, to)
		})
	}

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
