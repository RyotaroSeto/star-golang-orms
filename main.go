package main

import (
	"log"
	"star-golang-orms/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
