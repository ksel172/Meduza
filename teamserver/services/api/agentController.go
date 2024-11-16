package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/models"
)

type AgentController struct {
	dal *storage.AgentDAL
}

func NewAgentController(dal *storage.AgentDAL) *AgentController {
	return &AgentController{dal: dal}
}

func (ac *AgentController) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	var agent models.Agent

	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Set first contact variables
	agent.FirstCallback = time.Now()
	agent.ModifiedAt = time.Now()

	if err := ac.dal.RegisterAgent(agent); err != nil {
		http.Error(w, fmt.Sprintf("Error registering agent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return a response containing the agent for updating the client side
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agent)
}

func (ac *AgentController) GetAgent(w http.ResponseWriter, r *http.Request)    {}
func (ac *AgentController) UpdateAgent(w http.ResponseWriter, r *http.Request) {}
func (ac *AgentController) DeleteAgent(w http.ResponseWriter, r *http.Request) {}
