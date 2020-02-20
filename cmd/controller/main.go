package main

import (
	"log"

	"github.com/the-maldridge/yurt-tools/internal/controller"
)

func main() {
	c, err := controller.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Discover(); err != nil {
		log.Printf("Error during task discovery: %v", err)
	}
}
