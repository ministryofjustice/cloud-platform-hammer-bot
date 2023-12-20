package main

import (
	"log"

	"github.com/ministryofjustice/cloud-platform-hammer-bot/init_app"
)

func main() {
	ginMode, ghToken := init_app.InitEnvVars()

	ghClient, ghErr := init_app.InitGH(ghToken)
	if ghErr != nil {
		log.Fatal("Error initialising github client: ", ghErr)
	}

	ghRepo, ghErr := init_app.InitCommit()
	if ghErr != nil {
		log.Fatal("Error initialising github repo: ", ghErr)
	}

	r := init_app.InitGin(ginMode, ghClient, ghRepo, ghToken)

	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
