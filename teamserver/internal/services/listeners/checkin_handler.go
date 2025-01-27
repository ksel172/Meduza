package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ICheckInController interface {
	Checkin(ctx *gin.Context)
}

type CheckInController struct {
	checkInDAL dal.ICheckInDAL
	agentDAL   dal.IAgentDAL
}

func NewCheckInController(checkInDAL dal.ICheckInDAL, agentDAL dal.IAgentDAL) *CheckInController {
	return &CheckInController{checkInDAL: checkInDAL, agentDAL: agentDAL}
}

// need to protect by authentication at some points, because currently anyone requesting
// the tasks will get them, however, only the agent should be able to.
func (cc *CheckInController) Checkin(ctx *gin.Context) {

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c2request.Reason == models.Task {
		tasks, err := cc.agentDAL.GetAgentTasks(ctx, c2request.AgentID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		for i, task := range tasks {
			moduleDirPath := filepath.Join(conf.GetModuleUploadPath(), task.Module)
			moduleName := task.Command.Name

			modulePath := filepath.Join(moduleDirPath, moduleName)
			mainModuleBytes, err := utils.LoadAssembly(filepath.Join(modulePath, moduleName+".dll"))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load main module: %s", err.Error())})
				return
			}

			loadingModulePath := moduleDirPath + "/" + moduleName + "/"
			dependencyBytes := make(map[string][]byte)
			files, err := os.ReadDir(loadingModulePath)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read module directory: %s", err.Error())})
				return
			}

			for _, file := range files {
				if file.Name() != moduleName+".dll" && strings.HasSuffix(file.Name(), ".dll") {
					depBytes, err := utils.LoadAssembly(filepath.Join(loadingModulePath, file.Name()))
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load dependency: %s", err.Error())})
						return
					}
					dependencyBytes[file.Name()] = depBytes
				}
			}

			moduleBytes := models.ModuleBytes{
				ModuleBytes:     mainModuleBytes,
				DependencyBytes: dependencyBytes,
			}

			moduleBytesJSON, err := json.Marshal(moduleBytes)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to marshal module bytes: %s", err.Error())})
				return
			}

			task.Module = base64.StdEncoding.EncodeToString(moduleBytesJSON)
			tasks[i] = task
		}

		// Update the agent's last callback time
		lastCallback := time.Now().Format(time.RFC3339)
		if err := cc.agentDAL.UpdateAgentLastCallback(ctx.Request.Context(), c2request.AgentID, lastCallback); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal tasks to JSON"})
			return
		}

		var c2response models.C2Request
		c2response.AgentID = c2request.AgentID
		c2response.Reason = models.Task
		c2response.Message = string(tasksJSON)
		ctx.JSON(http.StatusOK, c2response)
	} else if c2request.Reason == models.Response {

		var agentTask models.AgentTask
		if jsonErr := json.Unmarshal([]byte(c2request.Message), &agentTask); jsonErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent info"})
			return
		}

		updateErr := cc.agentDAL.UpdateAgentTask(ctx.Request.Context(), agentTask)
		if updateErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
			return
		}

		ctx.JSON(http.StatusOK, "successfully updated")
	} else if c2request.Reason == models.Register {

		logger.Info("Received check-in request from agent:", c2request.AgentID)
		// Parse the message as AgentInfo
		var agentInfo models.AgentInfo
		if err := json.Unmarshal([]byte(c2request.Message), &agentInfo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent info"})
			return
		}
		// Check if the agent already exists
		if _, err := cc.agentDAL.GetAgent(agentInfo.AgentID); err == nil {
			logger.Info("Agent already exists:", c2request.AgentID)
			ctx.JSON(http.StatusConflict, gin.H{"error": "agent already exists"})
			return
		}

		// Create agent in the database
		newAgent := c2request.IntoNewAgent()
		newAgent.Name = utils.RandomString(6)

		if err := cc.checkInDAL.CreateAgent(ctx.Request.Context(), newAgent); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create agent info in the database
		if err := cc.agentDAL.CreateAgentInfo(ctx.Request.Context(), agentInfo); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"agent": newAgent})
	}
}
