package main

import (
	"log"

	"github.com/bluejutzu/GoBot/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
