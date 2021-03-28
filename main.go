package main

import (
	"log"

	"github.com/ecator/gofile/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
