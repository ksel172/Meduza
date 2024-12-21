package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type CheckInController struct {
	checkInDAL *dal.CheckInDAL
	agentDAL   *dal.AgentDAL
}

func NewCheckInController(checkInDAL *dal.CheckInDAL, agentDAL *dal.AgentDAL) *CheckInController {
	return &CheckInController{checkInDAL: checkInDAL, agentDAL: agentDAL}
}

func (cc *CheckInController) CreateAgent(ctx *gin.Context) {

	// Decode the received JSON into a C2Request
	// NewC2Request sets agentStatus as uninitialized if that is not provided by the agent in the JSON
	c2request := models.NewC2Request()
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if the received C2Request is valid
	if !c2request.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Convert C2Request into Agent model
	agent := c2request.IntoNewAgent()

	// Create agent in the redis db
	if err := cc.checkInDAL.CreateAgent(ctx.Request.Context(), agent); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"agent": agent})
}

// GetTasks Will be called by the agents to get their tasks/commands
// The agent will send its ID in the query params,
// need to protect by authentication at some points, because currently anyone requesting
// the tasks will get them, however, only the agent should be able to.
func (cc *CheckInController) GetTasks(ctx *gin.Context) {

	// Get the agent ID from the query params
	agentID := ctx.Param("id")
	if agentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "agent_id is required"})
		return
	}

	// Get the tasks for the agent
	tasks, err := cc.agentDAL.GetAgentTasks(ctx, agentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
