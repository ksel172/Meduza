package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/models"
)

type CheckInController struct {
	checkInDAL *redis.CheckInDAL
	agentDAL   *redis.AgentDAL
}

func NewCheckInController(checkInDAL *redis.CheckInDAL, agentDAL *redis.AgentDAL) *CheckInController {
	return &CheckInController{checkInDAL: checkInDAL, agentDAL: agentDAL}
}

func (cc *CheckInController) CreateAgent(w http.ResponseWriter, r *http.Request) {

	// Decode the received JSON into a C2Request
	// NewC2Request sets agentStatus as uninitialized if that is not provided by the agent in the JSON
	c2request := models.NewC2Request()
	if err := json.NewDecoder(r.Body).Decode(&c2request); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
	}

	// Validate if the received C2Request is valid
	if !c2request.Valid() {
		http.Error(w, "Invalid C2 request", http.StatusBadRequest)
		return
	}

	// Convert C2Request into Agent model
	agent := c2request.IntoNewAgent()

	// Create agent in the redis db
	if err := cc.checkInDAL.CreateAgent(agent); err != nil {
		http.Error(w, fmt.Sprintf("Error registering agent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return a response containing the agent for updating the client side
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agent)
}

// Will be called by the agents to get their tasks/commands
// The agent will send its ID in the query params,
// need to protect by authentication at some points, because currently anyone requesting
// the tasks will get them, however, only the agent should be able to.
func (cc *CheckInController) GetTasks(w http.ResponseWriter, r *http.Request) {

	// Get the agent ID from the query params
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		http.Error(w, "Missing agent ID", http.StatusBadRequest)
		return
	}

	// Get the tasks for the agent
	tasks, err := cc.agentDAL.GetAgentTasks(r.Context(), agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	// Return tasks as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
