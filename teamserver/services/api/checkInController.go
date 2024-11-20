package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	if err := cc.dal.CreateAgent(agent); err != nil {
		http.Error(w, fmt.Sprintf("Error registering agent: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return a response containing the agent for updating the client side
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agent)
}
