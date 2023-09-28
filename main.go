package main

import (
	"log"

	"github.com/ministryofjustice/cloud-platform-hammer-bot/init_app"
)

func main() {
	r := init_app.InitGin()

	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
