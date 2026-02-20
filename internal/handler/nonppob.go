package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/handler/model"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/internal/usecase"
	"git-rbi.jatismobile.com/databasemanagements/coster/api_ai_llm.git/utility"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type NonPPOBHandler interface {
	Create() gin.HandlerFunc
}

type NonPPOBHandlerImpl struct {
	nonPPOBUc usecase.NonPPOBUsecase
}

func NewNonPPOBHandler(nonPPOBUc usecase.NonPPOBUsecase) NonPPOBHandler {
	return &NonPPOBHandlerImpl{
		nonPPOBUc: nonPPOBUc,
	}
}

func (h *NonPPOBHandlerImpl) Create() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		jsonBody, err := utility.ReadJsonBodyRequest(c)
		if err != nil {
			utility.PrintConsole(fmt.Sprintf("%v", err.Error()), "error")
			intHttpStatus, status := utility.ErrorHttpStatus(err.Error())
			utility.FormatResponse(c, "json", intHttpStatus, gin.H{"code": intHttpStatus, "status": status, "message": err.Error()})
			return
		}

		objData := model.NonPPOBOrder{}

		err = json.Unmarshal([]byte(jsonBody), &objData)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"body":  jsonBody,
			}).Error("[NonPPOBHandler.Create] Failed to unmarshal request body")
			utility.FormatResponse(c, "json", http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "status": "error", "message": "Invalid request body format"})
			return
		}

		// Store the full request payload as received
		objData.MetadataRaw = jsonBody

		err = h.nonPPOBUc.Create(&objData)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"order": objData,
			}).Error("[NonPPOBHandler.Create] Failed to create non-ppob order")

			intHttpStatus, status := utility.ErrorHttpStatus(err.Error())
			utility.FormatResponse(c, "json", intHttpStatus, gin.H{"code": intHttpStatus, "status": status, "message": err.Error()})
			return
		}

		intHttpStatus := http.StatusOK
		utility.FormatResponse(c, "json", intHttpStatus, gin.H{"code": intHttpStatus, "message": "Success", "status": "success"})
	}
	return gin.HandlerFunc(fn)
}
