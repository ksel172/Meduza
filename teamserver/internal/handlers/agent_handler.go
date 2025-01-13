package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type AgentController struct {
	dal dal.IAgentDAL
}

func NewAgentController(dal dal.IAgentDAL) *AgentController {
	return &AgentController{
		dal: dal,
	}
}

func (ac *AgentController) GetAgent(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	agent, err := ac.dal.GetAgent(agentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agent)
}

func (ac *AgentController) UpdateAgent(ctx *gin.Context) {
	var agentUpdateRequest models.UpdateAgentRequest
	if err := ctx.ShouldBindJSON(&agentUpdateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentUpdateRequest.ModifiedAt = time.Now()

	updatedAgent, err := ac.dal.UpdateAgent(ctx, agentUpdateRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAgent)
}

func (ac *AgentController) DeleteAgent(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	if err := ac.dal.DeleteAgent(ctx, agentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

/* Agent Task API */

func (ac *AgentController) CreateAgentTask(ctx *gin.Context) {
	agentID := ctx.Param("id")
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	var agentTaskRequest models.AgentTaskRequest
	if err := ctx.ShouldBindJSON(&agentTaskRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentTask := agentTaskRequest.IntoAgentTask()
	agentTask.AgentID = agentID

	if err := ac.dal.CreateAgentTask(ctx, agentTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent task creation failed: %s", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, agentTask)
}

func (ac *AgentController) UpdateAgentTask(ctx *gin.Context) {
	agentID := ctx.Param("id")
	taskID := ctx.Param("task_id")
	if agentID == "" || taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "agent_id and task_id are required"})
		return
	}

	var agentTask models.AgentTask
	if err := ctx.ShouldBindJSON(&agentTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentTask.AgentID = agentID
	agentTask.TaskID = taskID

	if err := ac.dal.UpdateAgentTask(ctx, agentTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent task update failed: %s", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, agentTask)
}

func (ac *AgentController) GetAgentTasks(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	tasks, err := ac.dal.GetAgentTasks(ctx, agentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (ac *AgentController) DeleteAgentTasks(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	if err := ac.dal.DeleteAgentTasks(ctx, agentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Agent tasks deleted successfully"})
}

func (ac *AgentController) DeleteAgentTask(ctx *gin.Context) {
	agentID := ctx.Param("id")
	taskID := ctx.Param("task_id")
	if agentID == "" || taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "agent_id and task_id are required"})
		return
	}

	if err := ac.dal.DeleteAgentTask(ctx, agentID, taskID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Agent task deleted successfully"})
}

func (ac *AgentController) CreateAgentConfig(ctx *gin.Context) {
	var agentConfig models.AgentConfig
	if err := ctx.ShouldBindJSON(&agentConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.dal.CreateAgentConfig(ctx, agentConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agentConfig)
}

func (ac *AgentController) GetAgentConfig(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	agentConfig, err := ac.dal.GetAgentConfig(ctx, agentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agentConfig)
}

func (ac *AgentController) UpdateAgentConfig(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	var agentConfig models.AgentConfig
	if err := ctx.ShouldBindJSON(&agentConfig); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.dal.UpdateAgentConfig(ctx, agentID, agentConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agentConfig)
}

func (ac *AgentController) DeleteAgentConfig(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	if err := ac.dal.DeleteAgentConfig(ctx, agentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Agent config deleted successfully"})
}

func (ac *AgentController) CreateAgentInfo(ctx *gin.Context) {
	var agentInfo models.AgentInfo
	if err := ctx.ShouldBindJSON(&agentInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.dal.CreateAgentInfo(ctx, agentInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agentInfo)
}

func (ac *AgentController) UpdateAgentInfo(ctx *gin.Context) {
	var agentInfo models.AgentInfo
	if err := ctx.ShouldBindJSON(&agentInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.dal.UpdateAgentInfo(ctx, agentInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agentInfo)
}

func (ac *AgentController) GetAgentInfo(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	agentInfo, err := ac.dal.GetAgentInfo(ctx, agentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, agentInfo)
}

func (ac *AgentController) DeleteAgentInfo(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	if err := ac.dal.DeleteAgentInfo(ctx, agentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Agent info deleted successfully"})
}
