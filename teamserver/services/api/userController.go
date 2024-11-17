package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ksel172/Meduza/teamserver/internal/storage"
)

type UserController struct {
	dal *storage.UserDAL
}

func NewUserController(dal *storage.UserDAL) *UserController {
	return &UserController{dal: dal}
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := uc.dal.GetUsers(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response as JSON: %s", err.Error()), http.StatusInternalServerError)
	}
}
