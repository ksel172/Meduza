package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

	// agentID in the query params - could be just pure JSON
	agentID := r.URL.Query().Get("id")
	if agentID == "" {
		http.Error(w, "Missing agent ID", http.StatusBadRequest)
		return
	}

	// JSON modified agent
	var agent models.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
	}

	// Set agentID
	agent.ID = agentID

	if err := ac.dal.UpdateAgent(r.Context(), agent); err != nil {
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
	var agentTask models.AgentTask
	if err := json.NewDecoder(r.Body).Decode(&agentTask); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
	}

	// Generate	uuid
	agentTask.ID = uuid.New().String()

	if err := ac.dal.CreateAgentTask(r.Context(), agentTask); err != nil {
		http.Error(w, fmt.Sprintf("Agent not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agentTask)
}
func (ac *AgentController) GetAgentTasks(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("id")

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
