package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/models"
)

type AgentController struct {
	dal *redis.AgentDAL
}

/* Agents API */

func NewAgentController(dal *redis.AgentDAL) *AgentController {
	return &AgentController{dal: dal}
}

func (ac *AgentController) GetAgent(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		http.Error(w, "Missing agent ID", http.StatusBadRequest)
		return
	}

	agent, err := ac.dal.GetAgent(r.Context(), agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(agent)
}

// This is technically vulnerable because it allows the client to update any agent
// Also, allows any field to be edited if the request is hand-crafted
func (ac *AgentController) UpdateAgent(w http.ResponseWriter, r *http.Request) {

	// get the ID of the agent to be updated in the query params
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		http.Error(w, "Missing agent ID", http.StatusBadRequest)
		return
	}

	// Get the agent that is going to be updated in the db
	agent, err := ac.dal.GetAgent(r.Context(), agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	// Get the JSON for the fields that can be updated in the agent
	// This prevents unintented modifications by the client manipulating the request JSON
	var agentUpdateRequest models.UpdateAgentRequest
	if err = json.NewDecoder(r.Body).Decode(&agentUpdateRequest); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Update the agent fields
	updatedAgent := agentUpdateRequest.IntoAgent(agent)

	// Provide the updated agent to the data layer
	if err := ac.dal.UpdateAgent(r.Context(), updatedAgent); err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ac *AgentController) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		http.Error(w, "Missing agent ID", http.StatusBadRequest)
		return
	}

	if err := ac.dal.DeleteAgent(r.Context(), agentID); err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/* Agent Task API */

func (ac *AgentController) CreateAgentTask(w http.ResponseWriter, r *http.Request) {

	// Get the agentID from the query params
	agentID := r.URL.Query().Get("agent_id")
	if agentID == "" {
		http.Error(w, "Missing agent ID", http.StatusBadRequest)
		return
	}

	// Create agentTaskRequest model
	var agentTaskRequest models.AgentTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&agentTaskRequest); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
	}

	// Convert into AgentTask model with default fields, uuid generation,...
	agentTask := agentTaskRequest.IntoAgentTask()

	// Create the task for the agent in the db
	if err := ac.dal.CreateAgentTask(r.Context(), agentTask); err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agentTask)
}
func (ac *AgentController) GetAgentTasks(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")

	tasks, err := ac.dal.GetAgentTasks(r.Context(), agentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	// Return tasks as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// Deletes all tasks for a single agent
func (ac *AgentController) DeleteAgentTasks(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")

	if err := ac.dal.DeleteAgentTasks(r.Context(), agentID); err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete a single task
func (ac *AgentController) DeleteAgentTask(w http.ResponseWriter, r *http.Request) {
	agent_id := r.URL.Query().Get("agent_id")
	taskID := r.URL.Query().Get("task_id")

	if err := ac.dal.DeleteAgentTask(r.Context(), agent_id, taskID); err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
