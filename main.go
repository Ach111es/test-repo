package main

import (
	"fmt"
	"io"
	stdlog "log"

	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/routers"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/utility"
	"github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
)

func main() {
	// Silence log outputs to terminal; logging is handled via middleware/file
	logrus.SetOutput(io.Discard)
	stdlog.SetOutput(io.Discard)

	utility.PrintConsole("API Sync Data PCT - Bento service started", "info")
	utility.PrintConsole("Loading application configuration", "info")

	configuration, errConfig := utility.LoadApplicationConfiguration("")
	if errConfig != nil {
		logrus.WithFields(logrus.Fields{"error": errConfig}).Fatal("Failed to load app configuration")
	} else {

		utility.PrintConsole("Application configuration loaded successfully", "info")

		gin.SetMode(gin.ReleaseMode)

		utility.PrintConsole(fmt.Sprintf("Server running on port: %v", configuration.Http.HttpPort), "info")
		utility.PrintConsole(fmt.Sprintf("Queue Host: %v", configuration.Queue.Host), "info")
		utility.PrintConsole(fmt.Sprintf("Queue Name: %v", configuration.Queue.QueueName), "info")
		utility.PrintConsole(fmt.Sprintf("API Keys configured: %d", len(configuration.APIKeys)), "info")

		router := gin.New()
		routersInit := routers.InitRouter(configuration, router)

		utility.PrintConsole("API Sync Data PCT - Bento service initialized successfully", "info")

		errServer := routersInit.Run(":" + configuration.Http.HttpPort)
		if errServer != nil {
			utility.PrintConsole(fmt.Sprintf("%v", errServer.Error()), "error")
		}

	}
}
