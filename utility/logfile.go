package utility

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func WriteToLogFile(logFilePath, fileName, logMessage string) error {
	currentTime := time.Now().Format("2006-01-02")

	logFilePath = fmt.Sprintf("%s%s_%s.log", logFilePath, currentTime, fileName)

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		if _, err := os.Create(logFilePath); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write log message to log file
	if _, err := f.WriteString(logMessage); err != nil {
		return err
	}

	return nil
}

func BuildLogMessage(startTime time.Time, requestTimeInMs float64, c *gin.Context, writer *ResponseWriterWithCapture, requestBody []byte, responseHeaders http.Header, responseBody string) string {
	requestHeaderLog := formatLogEntry(startTime, requestTimeInMs, c.Request.Method, c.Request.URL.Path, writer.status, "Request Headers", c.Request.Header.Clone())
	requestBodyLog := formatLogEntry(startTime, requestTimeInMs, c.Request.Method, c.Request.URL.Path, writer.status, "Request Body", strings.NewReplacer("\n", "", "\r", "", " ", "").Replace(string(requestBody)))
	responseHeaderLog := formatLogEntry(startTime, requestTimeInMs, c.Request.Method, c.Request.URL.Path, writer.status, "Response Headers", responseHeaders)
	responseBodyLog := formatLogEntry(startTime, requestTimeInMs, c.Request.Method, c.Request.URL.Path, writer.status, "Response Body", responseBody)

	return fmt.Sprintf("%s%s%s%s", requestHeaderLog, requestBodyLog, responseHeaderLog, responseBodyLog)
}

func formatLogEntry(startTime time.Time, requestTimeInMs float64, method, path string, status int, sectionName string, data interface{}) string {
	return fmt.Sprintf("%s costerv3_api_pct | %.6fms | %-7s %s | [%d] %s: %v\n", startTime.Format("2006/01/02 15:04:05"), requestTimeInMs, method, path, status, sectionName, data)
}

type ResponseWriterWithCapture struct {
	gin.ResponseWriter
	Body   *bytes.Buffer
	status int
}

func (w *ResponseWriterWithCapture) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err == nil {
		w.Body.Write(b)
	}
	return n, err
}

func (w *ResponseWriterWithCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
