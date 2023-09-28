package init_app

import (
	"log"
	"os"
)

func InitEnvVars() (string, string) {
	githubTokenVal, githubTokenPresent := os.LookupEnv("GITHUB_TOKEN")
	if githubTokenVal == "" || !githubTokenPresent {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	ginMode := "debug"
	ginModeVal, ginModePresent := os.LookupEnv("GIN_MODE")
	if ginModeVal == "" || !ginModePresent {
		os.Setenv("GIN_MODE", ginMode)
		ginModeVal = ginMode
	}

	return ginModeVal, githubTokenVal
}
