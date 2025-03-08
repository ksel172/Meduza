package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
)

type AgentController struct {
	agentDal  dal.IAgentDAL
	moduleDal dal.IModuleDAL
}

func NewAgentController(agentDal dal.IAgentDAL, moduleDal dal.IModuleDAL) *AgentController {
	return &AgentController{
		agentDal:  agentDal,
		moduleDal: moduleDal,
	}
}

/* Agent API */

func (ac *AgentController) GetAgent(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	agent, err := ac.agentDal.GetAgent(ctx.Request.Context(), agentID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get agent", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent retrieved successfully", agent)
}

func (ac *AgentController) GetAgents(ctx *gin.Context) {
	agents, err := ac.agentDal.GetAgents(ctx)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get agents", err.Error())
		return
	}
	models.ResponseSuccess(ctx, http.StatusOK, "Agents retrieved successfully", agents)
}

func (ac *AgentController) UpdateAgent(ctx *gin.Context) {
	var agentUpdateRequest models.UpdateAgentRequest
	if err := ctx.ShouldBindJSON(&agentUpdateRequest); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	agentUpdateRequest.ModifiedAt = time.Now()

	updatedAgent, err := ac.agentDal.UpdateAgent(ctx, agentUpdateRequest)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update agent", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent updated successfully", updatedAgent)
}

func (ac *AgentController) DeleteAgent(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	if err := ac.agentDal.DeleteAgent(ctx, agentID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete agent", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent deleted successfully", nil)
}

/* Agent Task API */

func (ac *AgentController) CreateAgentTask(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	var agentTaskRequest models.AgentTaskRequest
	if err := ctx.ShouldBindJSON(&agentTaskRequest); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	agentTask := agentTaskRequest.IntoAgentTask()
	agentTask.AgentID = agentID

	if agentTask.Type == models.HelpCommand {

		agentTask.Started = time.Now()

		helpText := "Available commands:\n" +
			"shell [command] - Execute a shell command\n" +
			"help - brings up the help menu\n"

		modules, err := ac.moduleDal.GetAllModules(ctx)
		if err != nil {
			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get modules", err.Error())
			return
		}

		if modules != nil {
			helpText += "\nAvailable modules:\n"

			for _, module := range modules {

				helpText += fmt.Sprintf("\n%s\n", module.Name)

				for _, command := range module.Commands {
					helpText += fmt.Sprintf("- \t%s", command.Description)
				}

			}
		}

		agentTask.Status = models.TaskComplete
		agentTask.Finished = time.Now()
		agentTask.Command.Output = helpText
	}

	if err := ac.agentDal.CreateAgentTask(ctx, agentTask); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create agent task", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Agent task created successfully", agentTask)
}

func (ac *AgentController) UpdateAgentTask(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	taskID := ctx.Param(models.ParamTaskID)
	if agentID == "" || taskID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s and %s are required", models.ParamAgentID, models.ParamTaskID))
		return
	}

	var agentTask models.AgentTask
	if err := ctx.ShouldBindJSON(&agentTask); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	agentTask.AgentID = agentID
	agentTask.TaskID = taskID

	if err := ac.agentDal.UpdateAgentTask(ctx, agentTask); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update agent task", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent task updated successfully", agentTask)
}

func (ac *AgentController) GetAgentTasks(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	tasks, err := ac.agentDal.GetAgentTasks(ctx, agentID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get agent tasks", err.Error())
		return
	}

	// Load modules for each task
	for i, task := range tasks {
		//if task.Module != "" {
		modulePath := filepath.Join(conf.GetModuleUploadPath(), task.Module)

		moduleConfig, err := LoadModuleConfig(modulePath)
		if err != nil {
			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to load module configuration", err.Error())
			return
		}

		moduleBytes, err := json.Marshal(moduleConfig.Module)
		if err != nil {
			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to process module configuration", err.Error())
			return
		}

		task.Module = base64.StdEncoding.EncodeToString(moduleBytes)
		tasks[i] = task
		//}
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent tasks retrieved successfully", tasks)
}

func (ac *AgentController) DeleteAgentTasks(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	if err := ac.agentDal.DeleteAgentTasks(ctx, agentID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete agent tasks", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent tasks deleted successfully", nil)
}

func (ac *AgentController) DeleteAgentTask(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	taskID := ctx.Param(models.ParamTaskID)
	if agentID == "" || taskID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s and %s are required", models.ParamAgentID, models.ParamTaskID))
		return
	}

	if err := ac.agentDal.DeleteAgentTask(ctx, agentID, taskID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete agent task", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent task deleted successfully", nil)
}

func (ac *AgentController) CreateAgentConfig(ctx *gin.Context) {
	var agentConfig models.AgentConfig
	if err := ctx.ShouldBindJSON(&agentConfig); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := ac.agentDal.CreateAgentConfig(ctx, agentConfig); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create agent config", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Agent config created successfully", agentConfig)
}

func (ac *AgentController) GetAgentConfig(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	agentConfig, err := ac.agentDal.GetAgentConfig(ctx, agentID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get agent config", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent config retrieved successfully", agentConfig)
}

func (ac *AgentController) UpdateAgentConfig(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	var agentConfig models.AgentConfig
	if err := ctx.ShouldBindJSON(&agentConfig); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := ac.agentDal.UpdateAgentConfig(ctx, agentID, agentConfig); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update agent config", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent config updated successfully", agentConfig)
}

func (ac *AgentController) DeleteAgentConfig(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	if err := ac.agentDal.DeleteAgentConfig(ctx, agentID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete agent config", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent config deleted successfully", nil)
}

func (ac *AgentController) CreateAgentInfo(ctx *gin.Context) {
	var agentInfo models.AgentInfo
	if err := ctx.ShouldBindJSON(&agentInfo); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := ac.agentDal.CreateAgentInfo(ctx, agentInfo); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create agent info", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Agent info created successfully", agentInfo)
}

func (ac *AgentController) UpdateAgentInfo(ctx *gin.Context) {
	var agentInfo models.AgentInfo
	if err := ctx.ShouldBindJSON(&agentInfo); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := ac.agentDal.UpdateAgentInfo(ctx, agentInfo); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update agent info", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent info updated successfully", agentInfo)
}

func (ac *AgentController) GetAgentInfo(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	agentInfo, err := ac.agentDal.GetAgentInfo(ctx, agentID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get agent info", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent info retrieved successfully", agentInfo)
}

func (ac *AgentController) DeleteAgentInfo(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamAgentID))
		return
	}

	if err := ac.agentDal.DeleteAgentInfo(ctx, agentID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete agent info", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent info deleted successfully", nil)
}
