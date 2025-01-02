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
		ctx.JSON(http.StatusNotFound, fmt.Sprintf("Agent %s not found: %s", agentID, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, agent)
}

func (ac *AgentController) UpdateAgent(ctx *gin.Context) {

	// Get the JSON for the fields that can be updated in the agent
	// This prevents unintended modifications by the client manipulating the request JSON
	var agentUpdateRequest models.UpdateAgentRequest
	if err := ctx.ShouldBindJSON(&agentUpdateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update modified at
	agentUpdateRequest.ModifiedAt = time.Now()

	// Provide the updated agent to the data layer
	updatedAgent, err := ac.dal.UpdateAgent(ctx, agentUpdateRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent update failed: %s", err.Error()))
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
		ctx.JSON(http.StatusNotFound, fmt.Sprintf("Agent not found: %s", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

/* Agent Task API */

func (ac *AgentController) CreateAgentTask(ctx *gin.Context) {

	// Get the agentID from the query params
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	// Create agentTaskRequest model
	var agentTaskRequest models.AgentTaskRequest
	if err := ctx.ShouldBindJSON(&agentTaskRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert into AgentTask model with default fields, uuid generation,...
	agentTask := agentTaskRequest.IntoAgentTask()

	// Create the task for the agent in the db
	if err := ac.dal.CreateAgentTask(ctx, agentTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent task creation failed: %s", err.Error()))
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
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent task list failed: %s", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// DeleteAgentTasks Deletes all tasks for a single agent
func (ac *AgentController) DeleteAgentTasks(ctx *gin.Context) {
	agentID := ctx.Param(models.ParamAgentID)
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	if err := ac.dal.DeleteAgentTasks(ctx, agentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent task list failed: %s", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// DeleteAgentTask Delete a single task
func (ac *AgentController) DeleteAgentTask(ctx *gin.Context) {
	agentId := ctx.Param(models.ParamAgentID)
	if agentId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamAgentID)})
		return
	}

	taskID := ctx.Param(models.ParamTaskID)
	if taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", models.ParamTaskID)})
		return
	}

	if err := ac.dal.DeleteAgentTask(ctx, agentId, taskID); err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Agent task delete failed: %s", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Agent '%s' task '%s' deleted", agentId, taskID)})
}
