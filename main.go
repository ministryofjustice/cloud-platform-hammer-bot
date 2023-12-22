package main

import (
	"log"

	"github.com/ministryofjustice/cloud-platform-hammer-bot/init_app"
	"github.com/ministryofjustice/cloud-platform-hammer-bot/utils"
)

func main() {
	ginMode, ghToken, ghURL, ghUser := init_app.InitEnvVars()

	ghClient, ghErr := init_app.InitGH(ghToken)
	if ghErr != nil {
		log.Fatal("Error initialising github client: ", ghErr)
	}

	ghRepo, ghErr := init_app.InitCommit()
	if ghErr != nil {
		log.Fatal("Error initialising github repo: ", ghErr)
	}

	var gh = utils.GitHub{Mode: ginMode, Token: ghToken, URL: ghURL, User: ghUser, Repo: ghRepo, Client: ghClient}

	r := init_app.InitGin(gh)

	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
