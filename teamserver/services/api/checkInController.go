package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/models"
)

type CheckInController struct {
	dal *redis.CheckInDAL
}

func NewCheckInController(dal *redis.CheckInDAL) *CheckInController {
	return &CheckInController{dal: dal}
}

func (cc *CheckInController) CreateAgent(w http.ResponseWriter, r *http.Request) {
	var agent models.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request body: %s", err.Error()), http.StatusBadRequest)
	}

	// Validate required information
	if agent.Info.UUMOID == "" {
		http.Error(w, "Missing UUMOID", http.StatusBadRequest)
	}

	// Set first contact variables
	agent.FirstCallback = time.Now()
	agent.ModifiedAt = time.Now()

	// Create agent
	if err := cc.dal.CreateAgent(agent); err != nil {
		http.Error(w, fmt.Sprintf("Error registering agent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return a response containing the agent for updating the client side
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agent)
}
