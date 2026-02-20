package routers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/handler"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/handler/middleware"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/usecase"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/utility"
	"github.com/gin-gonic/gin"
)

// InitBentoRouter initializes the Bento API router with API Key authentication
func InitRouter(configuration utility.Configuration, router *gin.Engine) *gin.Engine {
	// Initialize usecases
	nonPPOBUc := usecase.NewNonPPOBUsecase(configuration)
	nonPPOBHandler := handler.NewNonPPOBHandler(nonPPOBUc)

	ppobUc := usecase.NewPPOBUsecase(configuration)
	ppobHandler := handler.NewPPOBHandler(ppobUc)

	// Logging middleware
	router.Use(func(c *gin.Context) {
		startTime := time.Now()

		requestBody, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		requestHeaders := c.Request.Header.Clone()
		buf := new(bytes.Buffer)
		writer := &bentoResponseWriterWithCapture{ResponseWriter: c.Writer, Body: buf}
		c.Writer = writer
		c.Next()

		latency := time.Since(startTime)
		requestTimeInMs := float64(latency.Nanoseconds()) / 1000000.0
		responseHeaders := c.Writer.Header()

		logMessage := fmt.Sprintf("%s %fms %-7s %s [%3d Request Headers: %v Request Body: %s Response Headers: %v Response Body: %s\n %v",
			startTime.Format("2006/01/02 15:04:05"),
			requestTimeInMs,
			c.Request.Method,
			c.Request.URL.Path,
			writer.status,
			requestHeaders,
			strings.NewReplacer("\n", "", "\r", "", " ", "").Replace(string(requestBody)),
			responseHeaders,
			buf.String(),
			c.Errors.String(),
		)

		var fileName string

		logFilePath := configuration.Log.LogFilePath
		if c.Writer.Status() == http.StatusOK {
			fileName = configuration.Log.DebugFileName
		} else {
			fileName = configuration.Log.ErrorFileName
		}

		err := utility.WriteToLogFile(logFilePath, fileName, logMessage)
		if err != nil {
			return
		}
	})

	gin.ForceConsoleColor()

	router.Use(utility.CORSMiddleware())
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "status": "error", "message": "Endpoint not found!"})
	})

	// Bento API Group with API Key authentication
	bentoApi := router.Group("/api/v1/orders")
	bentoApi.Use(middleware.APIKeyMiddleware(configuration.APIKeys))
	{
		bentoApi.POST("/non-ppob", nonPPOBHandler.Create())
		bentoApi.POST("/ppob", ppobHandler.Create())
	}

	return router
}

type bentoResponseWriterWithCapture struct {
	gin.ResponseWriter
	Body   *bytes.Buffer
	status int
}

func (w *bentoResponseWriterWithCapture) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err == nil {
		w.Body.Write(b)
	}
	return n, err
}

func (w *bentoResponseWriterWithCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
